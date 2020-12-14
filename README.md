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

- Remove images older than the date specified with option `--date`
- Keep images within a given retention period `--retention-period`
- Exclude some image repositories with option `--exclude-repository`
- Exclude images with certain tag(s) from deletion with option `--exclude-tag`
- Exclude images with tags matching a regexp pattern with option `--exclude-tag-pattern`
- Exclude images with tags matching a [SemVer](https://semver.org) pattern with option `--exclude-semver-tags`
  > Note: The SemVer standard does not include the `v` or `V` prefix (e.g. v1.0.0), but as it is widely used, our Regexp will also match tags beginning with either `v` or `V`, so they will be excluded from deletion as well.
- Only remove untagged images with `--untagged-only` flag
- Dry-run mode with option `--dry-run` (don't actually delete images but get same output)

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

```bash
gcrgc [options] <registry>
```

### Examples:

To cleanup the entire registry, run:
```bash
gcrgc gcr.io/project-id
```
> Warning ! All the repositories for that particular GCP project will be autodiscovered and cleaned.
Running this command will empty the entire registry !

You can set a retention period to keep the recent images and only clean the old ones.

Keep the images less than 30 days old:
```bash
gcrgc --retention-period=30d gcr.io/project-id
```

Can also be expressed with an absolute date:
```bash
gcrgc --date=2019-01-01 gcr.io/project-id
```

To limit the repositories to cleanup, you can either whitelist or blacklist a subset of repositories in the registry:

Cleanup the `gcr.io/project-id/nginx` and `gcr.io/project-id/my-app` repositories:
```bash
gcrgc --repositories=nginx,my-app gcr.io/project-id
```

Cleanup everything BUT the `gcr.io/project-id/nginx` and `gcr.io/project-id/my-app` repositories:
```bash
gcrgc --exclude-repositories=nginx,my-app gcr.io/project-id
```

You probably want to ensure the images with a certain tag are excluded from deletion:
```bash
gcrgc --exclude-tags=latest,other-tag gcr.io/project-id
```

Or, only clean untagged images:
```bash
gcrgc --untagged-only gcr.io/project-id
```

For more advanced control over tags exclution there are additional options:

Exclude tags matching a SemVer pattern (like `v1.0.0`):
```bash
gcrgc --exclude-semver-tags gcr.io/project-id
```

Exclude tags matching custom regexp patterns:
```bash
gcrgc \
  --exclude-tag-pattern '^release-.*' \
  --exclude-tag-pattern '^dev-.*' \
  gcr.io/project-id
```

### Using a configuration file

Instead of passing command-line flags, it is possible reference a configuration file instead:
```bash
gcrgc --config config.yaml
```

The config file matches the same structure as the command line options. Any option can be configured both in the command line and the configuration file.
The command line flags have a higher priority than the configuration defined in the file, so it's possible to override the file configuration with command line flags.

config.yaml:
```yaml
registry: gcr.io/project-id
retention-period: 30d
exclude-repositories:
  - nginx
  - my-app
exclude-semver-tags: true
exclude-tags:
  - latest
exclude-tag-pattern:
  - ^release-([0-9]+\.)+[0-9]+$
```

## Helm chart

A Helm chart is available if you wish to run it on a Kubernetes cluster (as a `CronJob`).

Check the [documentation](https://github.com/graillus/helm-charts/tree/master/charts/gcrgc)
