## gitstats
[![codecov](https://codecov.io/gh/stefanoj3/gitstats/branch/master/graph/badge.svg?token=VwDTLXudLv)](https://codecov.io/gh/stefanoj3/gitstats)

gitstats is a tool to collect statistics about github activity for a given set of: organization, repositories, users.

The application is still under heavy development, but the core features are already available.

## Table of Content
- [Getting started](#-getting-started)
    - [Download](#download)
    - [Running the software](#running-the-software)
    - [How to read the statistics](#how-to-read-the-statistics)
    - [How do I read the statistics](#how-do-i-read-the-statistics)
- [Development](#-development)
- [License](https://github.com/stefanoj3/gitstats/blob/master/LICENSE.md)

## [↑](#table-of-content) Getting Started

#### Download
You have 2 options to get the application:
- Docker: `docker run stefanoj3/gitstats:latest <cmd>`.
  `docker` needs to be available in your `$PATH`.
- Compile and manually install: fetch this repository, run `make build` and then `cp dist/gitstats <path-of-your-choice>`. 
  `go` and `make` have to be available in your `$PATH`.

#### Running the software
In order to run the software you first need to define a config file for your team.
[Here](sampleconfig.toml) you can find an example.

You also need to pass a valid github token to the application via the environment variable `GITHUB_TOKEN`.
The token must have the `repo` scope enabled.

If you have a config file, and the env variable is already available in your shell you can run the following command:
```bash
docker run -e GITHUB_TOKEN -v $PWD:/data stefanoj3/gitstats:latest gitstats collect -c /data/devconfig.toml -f 2020-04-01 -t 2020-05-01 -d 120h -v -o /data/out
```

If you want to look into how you can customize the execution you can run to get more info:
```bash
docker run -v $PWD:/data stefanoj3/gitstats:latest gitstats help
docker run -v $PWD:/data stefanoj3/gitstats:latest gitstats collect -h
```

The output of the command consists of 2 csv files:
- pull requests statistics for the time frame specified, and the users specified in the config file
- team statistics about pull requests open/closed/merged, and the average time to get a pr closed

You can then import the CSV files in your preferred sheet application and visualize the data. 

#### How to read the statistics
I use this tool to track trends in my team/projects.

This data alone do not represent the performance of the team or single individuals,
but you can use to track trends and knowing the context in which your team operates, 
better understand the whole picture. 

[Here](resources/docs/screenshot_users.png) you can see how the statistics for user can be visualized in a 
sheet application (using the *_user.csv output file).

To give you an idea hof how to read it, those are statistics about myself when I was onboarding on a new project,
as you can see there is a trend: the more time passes, the more I increased my contribution to the project.

I'd say this is a normal trend to observe when looking at the contribution of an engineer new to a project/organization.

[Here](resources/docs/out_pull_requests.csv) you can find an example output for pull requests statistics (the *_pull_requests.csv output file).
This is straight forward: how many PRs your team has opened(and still open), closed, and merged, in the given time frame.

It also tells you the average `TimeToMerge`: how much it takes for a PR to go from `open` to `merged` on average.

#### How do I read the statistics
As mentioned above, this tool help you to track trends, but does not reflect 1:1 
the performance of the people working on a given project.

**For example**, I track `TimeToMerge`, every month I build a report for the projects I'm interest into and 
look at this particular metric.

Is it increasing? Why? Can I do something to help my team get their PRs merged faster? 
What is the issue? Is CI too slow? Maybe I can look into it and find how to make it faster.

Is the project on an old codebase they know nothing about? Maybe I can find someone to mentor them.

Again, you need to know the context your team operates in to be able to consume this data.


## [↑](#table-of-content) Development

**Before Starting** run `make hook-install` to install a git hook that will help you to
avoid committing secrets by mistake.

Getting started is very easy, after you have cloned the repository you can launch `make help`
to see what commands are available and what do they do. 

You will need: `go`, `docker` and `make` available in your `$PATH` to e able to start developing.

A `GITHUB_TOKEN` env variable must be set with a token that has at least the `repo` scope.

In order to run gitstats while developing all you need to do is:
`GITHUBTOKEN=mytoken go run cmd/gitstats/main.go collect -c sampleconfig.toml -f 2020-01-01 -t 2020-01-31`

If you wanna play with the configuration just create another config file to use instead of `sampleconfig.toml`,
by default `devconfig.toml` is ignored in the `.gitignore` file. 
