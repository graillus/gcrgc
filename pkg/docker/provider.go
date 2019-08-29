package docker

// Provider can provide access to docker registry
type Provider interface {
	ListRepositories(registry string) []Repository
	ListImages(repo string, minDate string) []Image
	DeleteImage(repo string, img *Image, dryRun bool)
}
