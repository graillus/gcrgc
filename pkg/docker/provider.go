package docker

import "time"

// Provider can provide access to docker registry
type Provider interface {
	ListRepositories(registry string) []Repository
	ListImages(repo string, limit time.Time) []Image
	DeleteImage(repo string, img *Image, dryRun bool)
}
