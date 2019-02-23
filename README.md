gcrgc
=====

Clean up images on the Google Container Registry.

Initially based on the [gist](https://gist.github.com/ahmetb/7ce6d741bd5baa194a3fac6b1fec8bb7) by [Ahmet Alp Balkan](https://gist.github.com/ahmetb), and rewritten in Go.

## Prerequisites
authenticated `gcloud` session for the project.

## Usage

Clean up untagged images under the `gcr.io/project-id/my-image` repository.
```
gcrgc --repository=gcr.io/project-id/my-image --untagged-only
```

Clean up tagged and untagged images under the `gcr.io/project-id/my-image` repository older than 2019-01-01, excluding tags `master` and `latest`
```
gcrgc --repository=gcr.io/project-id/my-image --date=2019-01-01 --exclude-tag=latest --exclude-tag=master
```

Clean up tagged and untagged images under the entire project `gcr.io/project-id` older than 2019-01-01, excluding the images under `gcr.io/project-id/my-image`
```
gcrgc --project-registry=gcr.io/project-id --date=2019-01-01 --exclude-image=my-image
```
