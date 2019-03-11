package datasource

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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

	wg.Add(len(sources))
	for name, source := range sources {
		go func(name string, source Source) {
			defer wg.Done()
			if err := generateSource(dir, name, source); err != nil {
				mu.Lock()
				merr = multierror.Append(merr, err)
				mu.Unlock()
			}
		}(name, source)
	}
	wg.Wait()

	return merr.ErrorOrNil()
}

func generateSource(baseDir, name string, source Source) error {
	var wg sync.WaitGroup
	ch := make(chan interface{})

	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(baseDir, name))
	if err != nil {
		return err
	}
	defer out.Close()

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {
			enc.Encode(v)
		}
	}()
	err = source.Generate(ch)

	close(ch)
	wg.Wait()

	if err != nil {
		err = fmt.Errorf("%s: %s", name, err)
	}
	return err
}
