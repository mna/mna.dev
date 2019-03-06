# mna.dev

This is the source repository for my mna.dev personal website. It uses
a number of tools to generate data sources for relevant links about me
(e.g. twitter account, blog, articles, popular repositories, etc.) and
this gets turned by a build system into a web page that presents those
data sources as clickable "tiles" or "cards".

## list of data sources

* Twitter account
* 0value, hypermegatop blog websites
* Stack overflow profile
* Guest/job blog posts
  - https://blog.gopheracademy.com/advent-2014/goquery/
  - https://splice.com/blog/lesser-known-features-go-test/
  - https://splice.com/blog/going-extra-mile-golint-go-vet/
* Top-10 year-end music lists
* Horns of the devil website archive
* List of open source repos
  - github, gitlab, bitbucket, sr.ht
* Maybe the LinkedIn profile?
* Maybe short posts directly on this site?
* An about page

## features

* Super simple, straightforward: a single page with tiles for each entry, newest first
* Entries are generated by a build system (likely makefile, npm script or Go program)
* Filter by hashtags, search tiles
* About page
* Very little javascript, no cookies, no tracking ("analytics")

## dependencies

For Go:

* https://godoc.org/golang.org/x/oauth2
  - See for twitter: https://github.com/golang/oauth2/issues/175#issuecomment-299756224

Using this, it should be possible to retrieve data for:
* Github
* Gitlab
* Bitbucket
* StackOverflow
* Twitter (see linked comment)
* Possibly Sourcehut: https://man.sr.ht/meta.sr.ht/api.md

## install and generate

```
$ npm init
# installs development dependencies, populates node_modules

$ npm run build
# builds assets, generates css, js, html in public/

$ npm run serve
# starts a local web server to browse the static website, installs
# Go dependencies based on go.mod if needed

$ npm run generate
# runs the commands to retrieve data from the supported sources,
# updating existing ones as needed
```

Requires the following environments to be set (e.g. via an `.envrc` file
managed by `direnv`):

* GO111MODULE=on

