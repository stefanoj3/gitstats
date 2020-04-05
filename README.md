## gitstats
[![codecov](https://codecov.io/gh/stefanoj3/gitstats/branch/master/graph/badge.svg?token=VwDTLXudLv)](https://codecov.io/gh/stefanoj3/gitstats)

gitstats is a tool to collect statistics about github activity for a given: organization, repositories list, user handdles.

The application is still in development and it is missing core features.

## Table of Content
- [How to use it](#-how-to-use-it)
- [Download](#-download)
- [Development](#-development)
- [License](https://github.com/stefanoj3/gitstats/blob/master/LICENSE.md)

## [↑](#table-of-content) How to use it
TODO: 
- write me
- add example usage
- add screenshots


## [↑](#table-of-content) Download
TODO: write me

## [↑](#table-of-content) Development

**Before Starting** run `make hook-install` to install a git hook that will help you to
avoid committing secrets by mistake.

Getting started is very easy, after you have cloned the repository you can launch `make help`
to see what commands are available and what do they do. 

You will need: `go`, `docker` and `make` available in your `$PATH` to e able to start.

A `GITHUB_TOKEN` env variable must be set with a token that has at least the `repo` scope.

In order to run gitstats while developing all you need to do is:
`GITHUBTOKEN=mytoken go run cmd/gitstats/main.go collect -c sampleconfig.toml -f 2020-01-01 -t 2020-01-31`

If you wanna play with the configuration just create another config file to use instead of `sampleconfig.toml`,
by default `devconfig.toml` is ignored in the `.gitignore` file. 


## What is missing:
- the output needs to be printed in a format that can be used by a human
- no statistics per users are produced (only general statistics for all repos/users involved)
- add scrutinizer integration
- dockerize application & docker hub integration
