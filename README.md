gcrgc
=====

[![Build Status](https://travis-ci.org/graillus/gcrgc.svg?branch=master)](https://travis-ci.org/graillus/gcrgc)
[![codecov.io](http://codecov.io/github/graillus/gcrgc/coverage.svg?branch=master)](http://codecov.io/github/graillus/gcrgc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/graillus/gcrgc)](https://goreportcard.com/report/github.com/graillus/gcrgc)
[![GoDoc](https://godoc.org/github.com/graillus/gcrgc?status.svg)](https://godoc.org/github.com/graillus/gcrgc)
[![License MIT](https://img.shields.io/github/license/graillus/gcrgc.svg)](https://github.com/graillus/gcrgc/blob/master/LICENSE)

Clean up images on the Google Container Registry.

Initially based on the [gist](https://gist.github.com/ahmetb/7ce6d741bd5baa194a3fac6b1fec8bb7) by [Ahmet Alp Balkan](https://gist.github.com/ahmetb), and rewritten in Go.

## Features

- Remove images older than the date specified with option `-date`
- Clean up multiple image repositories at once with option `-all`
- Exclude some image repositories with option `-exclude-repository`
- Exclude images with certain tag(s) from deletion with option `-exclude-tag`
- Exclude images with tags matching a [SemVer](https://semver.org) pattern with option `-exclude-semver-tags`
  > Note: The SemVer standard does not include the `v` or `V` prefix (e.g. v1.0.0), but as it is widely used, our Regexp will also match tags beginning with either `v` or `V`, so they will be excluded from deletion as well.
- Only remove untagged images with `-untagged-only` flag
- Dry-run mode with option `-dry-run` (don't actually delete images but get same output)

## Prerequisites
authenticated `gcloud` session for the project.

## Installation

```
go get github.com/graillus/gcrgc
cd $GOPATH/src/github.com/graillus/gcrgc
make build
...
```

## Docker image

```
docker pull graillus/gcrgc
```

The docker image extends the google/cloud-sdk image, read the [documentation](https://hub.docker.com/r/google/cloud-sdk/) to learn how to authenticate using the docker image

## Usage

Clean up untagged images under the `gcr.io/project-id/my-image` repository.
```
gcrgc -registry=gcr.io/project-id -untagged-only my-image
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` repository older than 2019-01-01, excluding tags `master` and `latest`
```
gcrgc -registry=gcr.io/project-id -date=2019-01-01 -exclude-tag=latest -exclude-tag=master my-image
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` excluding SemVer tags and `latest`
```
gcrgc -registry=gcr.io/project-id -exclude-tag=latest -exclude-semver-tags my-image
```

Clean up tagged and untagged images under the entire registry `gcr.io/project-id` older than 2019-01-01, excluding the images under `gcr.io/project-id/my-image`
```
gcrgc -registry=gcr.io/project-id -all -date=2019-01-01 -exclude-repository=my-image
```
