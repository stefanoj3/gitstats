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

In order to run the software you first need to define a config file for your team.
[Here](sampleconfig.toml) you can find an example.

You also need to pass a valid github token to the application via the environment variable `GITHUB_TOKEN`.
The token must have the 

If you have a config file, and the env variable is already available in your shell you can run the following command:
```bash
docker run -e GITHUB_TOKEN -v $PWD:/data gitstats:latest gitstats collect -c /data/devconfig.toml -f 2020-04-01 -t 2020-05-01 -d 120h -v -o /data/out
```

If you want to look into how you can customize the execution you can run to get more info:
```bash
docker run -v $PWD:/data gitstats:latest gitstats help
docker run -v $PWD:/data gitstats:latest gitstats collect -h
```

The output of the command consists of 2 csv files:
- pull requests statistics for the time frame specified, and the users specified in the config file
- team statistics about pull requests open/closed/merged, and the average time to get a pr closed

You can then import the CSV files in your preferred sheet application and visualize the data. 

## [↑](#table-of-content) Download
You have 2 options to get the application:
- `docker run gitstats:latest <cmd>`: requires `docker` to be available in your `$PATH`
- Fetch this repository, run `make build` and then `cp dist/gitstats <path-of-your-choice>`: requires `go` and `make` to be available in your `$PATH`


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
