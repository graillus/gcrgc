gcrgc
=====

[![Build Status](https://travis-ci.org/graillus/gcrgc.svg?branch=master)](https://travis-ci.org/graillus/gcrgc)
[![Go Report Card](https://goreportcard.com/badge/github.com/graillus/gcrgc)](https://goreportcard.com/report/github.com/graillus/gcrgc)
[![GoDoc](https://godoc.org/github.com/graillus/gcrgc?status.svg)](https://godoc.org/github.com/graillus/gcrgc/cmd/gcrgc)
[![License MIT](https://img.shields.io/github/license/graillus/gcrgc.svg)](https://github.com/graillus/gcrgc/blob/master/LICENSE)

Clean up images on the Google Container Registry.

Initially based on the [gist](https://gist.github.com/ahmetb/7ce6d741bd5baa194a3fac6b1fec8bb7) by [Ahmet Alp Balkan](https://gist.github.com/ahmetb), and rewritten in Go.

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
gcrgc -repository=gcr.io/project-id -untagged-only my-image
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` repository older than 2019-01-01, excluding tags `master` and `latest`
```
gcrgc -repository=gcr.io/project-id -date=2019-01-01 -exclude-tag=latest -exclude-tag=master my-image
```

Clean up tagged and untagged images under the entire project `gcr.io/project-id` older than 2019-01-01, excluding the images under `gcr.io/project-id/my-image`
```
gcrgc -repository=gcr.io/project-id -all -date=2019-01-01 -exclude-image=my-image
```
