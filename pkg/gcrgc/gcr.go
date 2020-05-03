package gcrgc

import (
	"fmt"
	"log"
	"time"

	"github.com/graillus/gcrgc/pkg/docker"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GCR abstracts the google/go-containerregistry SDK
type GCR struct {
	auth authn.Authenticator
}

// NewGCR creates a new instance of GCloud
func NewGCR(auth authn.Authenticator) *GCR {
	return &GCR{auth}
}

// ListRepositories gets the list of repositories for current registry
func (g GCR) ListRepositories(registry string) []docker.Repository {
	repo, err := name.NewRepository(registry)
	if err != nil {
		log.Fatalf("Cannot create repository: %s\n", err)
	}

	tags, err := google.List(repo, google.WithAuth(g.auth))
	if err != nil {
		log.Fatalf("Cannot list repository: %s\n", err)
	}

	var repos []docker.Repository
	for _, item := range tags.Children {
		repos = append(repos, *docker.NewRepository(registry + "/" + item))
	}

	return repos
}

// ListImages gets the list of images for the given repository name
func (g GCR) ListImages(reponame string, limit time.Time) []docker.Image {
	repo, err := name.NewRepository(reponame)
	if err != nil {
		log.Fatalf("Cannot create repository: %s\n", err)
	}

	tags, err := google.List(repo, google.WithAuth(g.auth))
	if err != nil {
		log.Fatalf("Cannot list repository: %s\n", err)
	}

	var images []docker.Image
	for hash, manifest := range tags.Manifests {
		if manifest.Uploaded.UTC().Before(limit) {
			images = append(images, *docker.NewImage(hash, manifest.Tags))
		}
	}

	return images
}

// DeleteImage deletes an image
func (g GCR) DeleteImage(repo string, i *docker.Image, dryRun bool) {
	if dryRun {
		i.IsRemoved = true

		return
	}

	for _, tag := range i.Tags {
		err := g.deleteItem(repo + ":" + tag)
		if err != nil {
			return
		}
	}

	err := g.deleteItem(repo + "@" + i.Digest)
	if err != nil {
		return
	}

	i.IsRemoved = true
}

func (g GCR) deleteItem(item string) error {
	ref, err := name.ParseReference(item)
	if err != nil {
		return fmt.Errorf("Unable to parse reference %s: %s", ref, err)
	}

	err = remote.Delete(ref, remote.WithAuth(g.auth))
	if err != nil {
		return fmt.Errorf("Unable to delete tag %s: %s", ref, err)
	}

	return nil
}
