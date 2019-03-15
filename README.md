# mna.dev

This is the source repository for my mna.dev personal website. It uses
a number of tools to generate data sources for relevant links about me
(e.g. twitter account, blog, articles, popular repositories, etc.) and
this gets turned by a build system into a web page that presents those
data sources as clickable "tiles" or "cards".

## TODOs

* Filter by hashtags
* Link directly to hashtag filters
* Search in twitter cards
* Grab image(s) when extracting the data, render more twitter-like cards

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

* `GO111MODULE=on` (until this is the default)
* `GITHUB_API_TOKEN`
* `SRHT_API_TOKEN`
* `GITLAB_API_TOKEN`
* `BITBUCKET_API_TOKEN`

