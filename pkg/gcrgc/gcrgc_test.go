package gcrgc

import (
	"testing"
	"time"

	"github.com/graillus/gcrgc/pkg/docker"
)

type fakeProvider struct{}

func (p fakeProvider) ListRepositories(registry string) []docker.Repository {
	return []docker.Repository{}
}

func (p fakeProvider) ListImages(repo string, limit time.Time) []docker.Image {
	return []docker.Image{
		*docker.NewImage("image", []string{"tag"}),
		*docker.NewImage("image-excluded-semver-tag", []string{"1.2.3"}),
		*docker.NewImage("image-excluded-tag-pattern", []string{"excluded-pattern-01"}),
		*docker.NewImage("image-excluded-tag", []string{"excluded"}),
	}
}

func (p fakeProvider) DeleteImage(repo string, img *docker.Image, dryRun bool) {}

func TestGetsTaskList(t *testing.T) {
	var settings = Settings{
		Registry:            "gcr.io/foo",
		AllRepositories:     true,
		Date:                nil,
		ExcludeSemVerTags:   true,
		ExcludedTagPatterns: []string{"^excluded-pattern-[0-9]*$"},
		ExcludedTags:        []string{"excluded"},
	}
	var (
		provider     = &fakeProvider{}
		repositories = []docker.Repository{*docker.NewRepository("image")}
	)

	tasks := getTaskList(provider, repositories, &settings)
	if len(tasks) != 1 {
		t.Errorf("Expected tasks number to be %d, got %d instead", 1, len(tasks))
	}

	if _, ok := tasks["image"]; ok == false {
		t.Errorf("Expected task key to be the name of the repository")
	}
}
