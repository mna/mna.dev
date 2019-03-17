# mna.dev

This is the source repository for my mna.dev personal website. It uses
a number of tools to generate data sources for relevant links about me
(e.g. twitter account, blog, articles, popular repositories, etc.) and
this gets turned by a build system into a web page that presents those
data sources as clickable "tiles" or "cards".

## TODOs

* Grab image(s) when extracting the data, render more twitter-like cards
* Click anywhere-ish on a card should link to the URL?
* Un-rendered twitter embeds should have a Twitter indication/card-header?

## install and generate

```
$ npm install
# installs development dependencies, populates node_modules

$ npm run build
# cleans output, generates data, builds assets, generates css, js, html in public/

$ npm run serve
# starts a local web server to browse the static website, installs
# Go dependencies based on go.mod if needed

$ npm run generate
# runs the commands to retrieve data from the supported sources,
# updating existing ones as needed

$ npm-run-all build:*
# only build the website, do not re-generate data and images

$ npm run watch
# watch for changes and run corresponding build task live
```

Requires the following environments to be set (e.g. via an `.envrc` file
managed by `direnv`):

* `GO111MODULE=on` (until this is the default)
* `GITHUB_API_TOKEN`
* `SRHT_API_TOKEN`
* `GITLAB_API_TOKEN`
* `BITBUCKET_API_TOKEN`

