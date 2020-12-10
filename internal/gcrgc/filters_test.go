package gcrgc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/graillus/gcrgc/internal/docker"
)

func TestGetRepoList(t *testing.T) {
	t.Run("One repository included", func(t *testing.T) {
		s := &Settings{
			AllRepositories: false,
			Repositories:    []string{"repo1"},
		}

		expected := []docker.Repository{
			docker.Repository{Name: "gcr.io/project/repo1"},
		}
		actual := getRepositoryList(createFakeRegistry(), s)

		err := assertRepositories(expected, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("One unexsisting repository included", func(t *testing.T) {
		s := &Settings{
			AllRepositories: false,
			Repositories:    []string{"unknown"},
		}

		expected := []docker.Repository{}
		actual := getRepositoryList(createFakeRegistry(), s)

		err := assertRepositories(expected, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("All repositories included, none excluded", func(t *testing.T) {
		s := &Settings{AllRepositories: true}

		expected := []docker.Repository{
			docker.Repository{Name: "gcr.io/project/repo1"},
			docker.Repository{Name: "gcr.io/project/repo2"},
			docker.Repository{Name: "gcr.io/project/repo3"},
		}
		actual := getRepositoryList(createFakeRegistry(), s)

		err := assertRepositories(expected, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("All repositories included, one excluded", func(t *testing.T) {
		s := &Settings{
			AllRepositories:      true,
			ExcludedRepositories: []string{"repo2"},
		}

		expected := []docker.Repository{
			docker.Repository{Name: "gcr.io/project/repo1"},
			docker.Repository{Name: "gcr.io/project/repo3"},
		}
		actual := getRepositoryList(createFakeRegistry(), s)

		err := assertRepositories(expected, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("All repositories included, one unexisting excluded", func(t *testing.T) {
		s := &Settings{
			AllRepositories:      true,
			ExcludedRepositories: []string{"unknown"},
		}

		expected := []docker.Repository{
			docker.Repository{Name: "gcr.io/project/repo1"},
			docker.Repository{Name: "gcr.io/project/repo2"},
			docker.Repository{Name: "gcr.io/project/repo3"},
		}
		actual := getRepositoryList(createFakeRegistry(), s)

		err := assertRepositories(expected, actual)
		if err != nil {
			t.Error(err)
		}
	})
}

func createFakeRegistry() *docker.Registry {
	repos := []docker.Repository{
		docker.Repository{Name: "gcr.io/project/repo1"},
		docker.Repository{Name: "gcr.io/project/repo2"},
		docker.Repository{Name: "gcr.io/project/repo3"},
	}

	return docker.NewRegistry("gcr.io/project", repos)
}

func assertRepositories(expected []docker.Repository, actual []docker.Repository) error {
	if len(expected) != len(actual) {
		msg := fmt.Sprintf("Expected repositories to contain %d elements, got %d instead", len(expected), len(actual))

		return errors.New(msg)
	}

	for _, e := range expected {
		found := false
		for _, a := range actual {
			if e.Name == a.Name {
				found = true
			}
		}

		if found == false {
			msg := fmt.Sprintf("Repositories does not contain element %s", e.Name)

			return errors.New(msg)
		}
	}

	return nil
}

func TestFilterImages(t *testing.T) {
	t.Run("No filter", func(t *testing.T) {
		expectedDigests := []string{"untagged", "foo", "foo-bar-baz"}

		filters := []ImageFilter{}

		actual := filterImages(createFakeImages(), filters)

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Untagged Only", func(t *testing.T) {
		expectedDigests := []string{"untagged"}

		filters := []ImageFilter{
			UntaggedFilter{true},
		}

		actual := filterImages(createFakeImages(), filters)

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Exclude tag", func(t *testing.T) {
		expectedDigests := []string{"untagged", "foo"}

		filters := []ImageFilter{
			UntaggedFilter{false},
			TagNameFilter{[]string{"bar"}},
		}

		actual := filterImages(createFakeImages(), filters)

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Exclude tag pattern", func(t *testing.T) {
		expectedDigests := []string{"untagged"}

		filters := []ImageFilter{
			UntaggedFilter{false},
			TagNameFilter{[]string{}},
			NewTagNameRegexFilter([]string{"^foo"}),
		}

		actual := filterImages(createFakeImages(), filters)

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Exclude SemVer tag pattern", func(t *testing.T) {
		fakeImages := []docker.Image{
			// Should match
			docker.Image{Digest: "invalid", Tags: []string{"2020-02-02-12345"}},
			docker.Image{Digest: "invalid", Tags: []string{"latest"}},
			docker.Image{Digest: "invalid", Tags: []string{"V1.0.0"}},
			// Should not match
			docker.Image{Digest: "0.0.0", Tags: []string{"0.0.0"}},
			docker.Image{Digest: "v1.0.0", Tags: []string{"v1.0.0"}},
			docker.Image{Digest: "V1.0.0", Tags: []string{"V1.0.0"}},
			docker.Image{Digest: "999.999.999", Tags: []string{"999.999.999"}},
			docker.Image{Digest: "v0.10", Tags: []string{"v0.10"}},
			docker.Image{Digest: "2.0.0-rc3", Tags: []string{"2.0.0-rc3"}},
		}

		expectedDigests := []string{
			"invalid",
			"invalid",
			"invalid",
		}

		filters := []ImageFilter{
			NewSemVerTagNameFilter(true),
		}

		actual := filterImages(fakeImages, filters)

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})
}

func createFakeImages() []docker.Image {
	return []docker.Image{
		docker.Image{Digest: "untagged", Tags: []string{}},
		docker.Image{Digest: "foo", Tags: []string{"foo"}},
		docker.Image{Digest: "foo-bar-baz", Tags: []string{"foo", "bar", "baz"}},
	}
}

func assertImages(expectedDigests []string, actual []docker.Image) error {
	if len(expectedDigests) != len(actual) {
		msg := fmt.Sprintf("Expected images to contain %d elements, got %d instead", len(expectedDigests), len(actual))

		return errors.New(msg)
	}

	for _, e := range expectedDigests {
		found := false
		for _, a := range actual {
			if e == a.Digest {
				found = true
			}
		}

		if found == false {
			msg := fmt.Sprintf("Images does not contain element with digest %s", e)

			return errors.New(msg)
		}
	}

	return nil
}
