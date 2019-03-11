package datasource

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
	"github.com/hashicorp/go-multierror"
)

var sources = make(map[string]Source)

// Source represents a data source that emits data values on
// the emit channel.
type Source interface {
	Generate(emit chan<- interface{}) error
}

// Register registers s as a data source under the specified name.
// It panics if a source is already registered for that name.
func Register(name string, source Source) {
	if _, ok := sources[name]; ok {
		panic(fmt.Errorf("a data source is already registered for %s", name))
	}
	sources[name] = source
}

// Generate calls Generate for each registered source, in parallel,
// and returns any error it encountered.
func Generate(dir string) error {
	var mu sync.Mutex
	var merr *multierror.Error
	var wg sync.WaitGroup

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	// create files for posts, micro-posts and repos
	writers := map[string]io.Writer{"post": nil, "mpost": nil, "repo": nil}
	for nm := range writers {
		f, err := os.Create(filepath.Join(dir, nm+".json"))
		if err != nil {
			return err
		}
		defer f.Close()

		bw := bufio.NewWriter(f)
		writers[nm] = &lockWriter{w: bw}
		defer bw.Flush()
	}

	wg.Add(len(sources))
	for name, source := range sources {
		go func(name string, source Source) {
			defer wg.Done()

			seri := &serializer{
				sourceName: name,
				source:     source,
				wpost:      writers["post"],
				wmpost:     writers["mpost"],
				wrepo:      writers["repo"],
			}
			if err := seri.serialize(); err != nil {
				mu.Lock()
				merr = multierror.Append(merr, err)
				mu.Unlock()
			}
		}(name, source)
	}
	wg.Wait()

	return merr.ErrorOrNil()
}

type lockWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (w *lockWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	n, err := w.w.Write(b)
	w.mu.Unlock()
	return n, err
}

type serializer struct {
	sourceName string
	source     Source
	encodeErr  error

	wpost  io.Writer
	wmpost io.Writer
	wrepo  io.Writer
}

func (s *serializer) serialize() error {
	var wg sync.WaitGroup
	ch := make(chan interface{})

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range ch {
			if s.encodeErr != nil {
				continue // drain ch
			}

			b, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				s.encodeErr = err
				continue
			}

			var werr error
			switch v.(type) {
			case *types.Post:
				_, werr = s.wpost.Write(b)
			case *types.MicroPost:
				_, werr = s.wmpost.Write(b)
			case *types.Repo:
				_, werr = s.wrepo.Write(b)
			}
			if werr != nil {
				s.encodeErr = werr
				continue
			}
		}
	}()
	err := s.source.Generate(ch)

	close(ch)
	wg.Wait()

	if err == nil {
		err = s.encodeErr
	}
	if err != nil {
		err = fmt.Errorf("%s: %s", s.sourceName, err)
	}
	return err
}
