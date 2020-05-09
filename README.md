gcrgc
=====

[![Build Status](https://travis-ci.org/graillus/gcrgc.svg?branch=master)](https://travis-ci.org/graillus/gcrgc)
[![codecov.io](http://codecov.io/github/graillus/gcrgc/coverage.svg?branch=master)](http://codecov.io/github/graillus/gcrgc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/graillus/gcrgc)](https://goreportcard.com/report/github.com/graillus/gcrgc)
[![GoDoc](https://godoc.org/github.com/graillus/gcrgc?status.svg)](https://godoc.org/github.com/graillus/gcrgc)
[![License MIT](https://img.shields.io/github/license/graillus/gcrgc.svg)](https://github.com/graillus/gcrgc/blob/master/LICENSE)

The GCR Garbage Collector

Tool for cleaning up images on the Google Container Registry.
Initially based on the [gist](https://gist.github.com/ahmetb/7ce6d741bd5baa194a3fac6b1fec8bb7) by [Ahmet Alp Balkan](https://gist.github.com/ahmetb), and rewritten in Go.

## Features

- Remove images older than the date specified with option `-date`
- Keep images within a given retention period `-retention-period`
- Clean up multiple image repositories at once with option `-all`
- Exclude some image repositories with option `-exclude-repository`
- Exclude images with certain tag(s) from deletion with option `-exclude-tag`
- Exclude images with tags matching a regexp pattern with option `-exclude-tag-pattern`
- Exclude images with tags matching a [SemVer](https://semver.org) pattern with option `-exclude-semver-tags`
  > Note: The SemVer standard does not include the `v` or `V` prefix (e.g. v1.0.0), but as it is widely used, our Regexp will also match tags beginning with either `v` or `V`, so they will be excluded from deletion as well.
- Only remove untagged images with `-untagged-only` flag
- Dry-run mode with option `-dry-run` (don't actually delete images but get same output)

## Prerequisites

You need an authenticated local `gcloud` installation, and write access to a Google Container Registry.

You can use a service account as well by setting the `GOOGLE_APPLICATION_CREDENTIALS` environment variable. Read the Google [documentation](https://cloud.google.com/docs/authentication/getting-started) for more details.

## Installation

### Binary releases

1. Download your [desired version](https://github.com/graillus/gcrgc/releases)
2. Extract it
```bash
tar xvf gcrgc_0.3.2_linux_amd64.tar.gz
```
3. Move binary to desired destination
```bash
mv gcrgc /usr/local/bin
```

### From sources

```bash
go get github.com/graillus/gcrgc
cd $GOPATH/src/github.com/graillus/gcrgc
go build -o bin/gcrgc cmd/gcrgc/gcrgc.go
```

### Using docker

A public image repository is available on [DockerHub](https://hub.docker.com/r/graillus/gcrgc)

```bash
docker pull graillus/gcrgc
```

Run with Google service account credentials:
```bash
docker run -t --rm \
  -v /path/to/serviceaccount.json:/credentials \
  -e GOOGLE_APPLICATION_CREDENTIALS=/credentials/serviceaccount.json
  graillus/gcrgc ...
```

## Usage

Clean up untagged images under the `gcr.io/project-id/my-image` repository.
```bash
gcrgc -registry=gcr.io/project-id -untagged-only my-image
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` repository older than 2019-01-01, excluding tags `master` and `latest`
```bash
gcrgc -registry=gcr.io/project-id -date=2019-01-01 -exclude-tag=latest -exclude-tag=master my-image
```

Clean up images older than 30 days
```bash
gcrgc -registry=gcr.io/project-id -retention-period 30d
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` excluding SemVer tags and `latest`
```bash
gcrgc -registry=gcr.io/project-id -exclude-tag=latest -exclude-semver-tags my-image
```

Clean up tagged and untagged images under the entire registry `gcr.io/project-id` older than 2019-01-01, excluding the images under `gcr.io/project-id/my-image`
```bash
gcrgc -registry=gcr.io/project-id -all -date=2019-01-01 -exclude-repository=my-image
```

## Helm chart

A Helm chart is available if you wish to run it on a Kubernetes cluster (as a `CronJob`).

Check the [documentation](https://github.com/graillus/helm-charts/tree/master/charts/gcrgc)
