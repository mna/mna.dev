package types

import (
	"html/template"
	"sort"
	"strings"
	"time"
)

// MarkdownPost is a post extracted from the .toml metadata file and
// the optional .md markdown file. In other words, it represents a post
// in-memory (in this struct), but not yet generated to a public/ page,
// so it doesn't have an URL yet.
type MarkdownPost struct {
	Path      string
	Title     string
	Published time.Time
	Lead      string
	Micro     bool
	Markdown  []byte
}

// ToPost converts a MarkdownPost to a Post.
func (mp *MarkdownPost) ToPost() *Post {
	p := &Post{
		URL:       "/" + mp.Path,
		Website:   "mna.dev",
		Title:     mp.Title,
		Lead:      mp.Lead,
		Published: mp.Published,
	}
	p.SetTags()
	return p
}

// ToMicroPost converts a MarkdownPost to a Post.
func (mp *MarkdownPost) ToMicroPost() *MicroPost {
	p := &MicroPost{
		Website:   "mna.dev",
		Text:      string(mp.Markdown),
		Published: mp.Published,
	}
	p.SetTags()
	return p
}

// PostConfig holds the configuration of a post as read from
// the toml files.
type PostConfig struct {
	Title     string    `toml:"title"`
	Lead      string    `toml:"lead"`
	Micro     bool      `toml:"micro"`
	Published time.Time `toml:"published"`
}

// Repo is the struct for a repository.
type Repo struct {
	// URL is the link to the repository.
	URL string
	// Host is the name of the host where it is hosted.
	Host string
	// Name is the name of the repository.
	Name string
	// Description is a short description of the repository.
	Description string
	// Language is the main programming language used in the
	// repository.
	Language string
	// Created is the date it was created.
	Created time.Time
	// Updated is the date it was last updated.
	Updated time.Time
	// Stars is the number of "stars" or "likes" of the repository
	// on the host website.
	Stars int
	// Forks is the number of "forks" of the repository on the
	// host website.
	Forks int
	// Tags is the list of tags associated with the repository.
	Tags []string
}

// SetTags sets the tags on r, adding default tags in
// addition to the provided tags.
func (r *Repo) SetTags(tags ...string) {
	r.Tags = append(r.Tags, tags...)
	r.Tags = append(r.Tags, "code", r.Host, r.Language)
	r.Tags = canonicalizeTags(r.Tags)
}

// Post is the struct for a blog post.
type Post struct {
	// URL is the link to the post.
	URL string
	// Website is the name of the website where this is hosted.
	Website string
	// Title is the title of the post.
	Title string
	// Lead is the short introduction of the post.
	Lead string
	// Published is the date it was published.
	Published time.Time
	// Tags is the list of tags associated with the post.
	Tags []string
}

// SetTags sets the tags on p, adding default tags in
// addition to the provided tags.
func (p *Post) SetTags(tags ...string) {
	p.Tags = append(p.Tags, tags...)
	p.Tags = append(p.Tags, "post", p.Website)
	p.Tags = canonicalizeTags(p.Tags)
}

// MicroPost is the struct for a micro-post.
type MicroPost struct {
	// URL is the link to that micro-post, if hosted elsewhere.
	URL string
	// Website is the name of the website where this is hosted.
	Website string
	// Text is the text-only content of the micro post.
	Text string
	// RawHTML contains the html markup to render this micro-post.
	RawHTML template.HTML
	// Published is the date it was published.
	Published time.Time
	// Tags is the list of tags associated with the micro post.
	Tags []string
}

// SetTags sets the tags on p, adding default tags in
// addition to the provided tags.
func (p *MicroPost) SetTags(tags ...string) {
	p.Tags = append(p.Tags, tags...)
	p.Tags = append(p.Tags, "micro", p.Website)
	p.Tags = canonicalizeTags(p.Tags)
}

func canonicalizeTags(tags []string) []string {
	set := make(map[string]bool)
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		set[t] = true
	}

	canon := make([]string, 0, len(set))
	for t := range set {
		canon = append(canon, t)
	}
	sort.Strings(canon)
	return canon
}
