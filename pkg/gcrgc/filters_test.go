package gcrgc

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/graillus/gcrgc/pkg/docker"
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
		actual := getRepoList(createFakeRegistry(), s)

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
		actual := getRepoList(createFakeRegistry(), s)

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
		actual := getRepoList(createFakeRegistry(), s)

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
		actual := getRepoList(createFakeRegistry(), s)

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
		actual := getRepoList(createFakeRegistry(), s)

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

func TestGetImageList(t *testing.T) {
	t.Run("No filter", func(t *testing.T) {
		expectedDigests := []string{"untagged", "foo", "foo-bar-baz"}

		actual := getImageList(createFakeImages(), false, []string{}, []*regexp.Regexp{})

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Untagged Only", func(t *testing.T) {
		expectedDigests := []string{"untagged"}

		actual := getImageList(createFakeImages(), true, []string{}, []*regexp.Regexp{})

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Exclude tag", func(t *testing.T) {
		expectedDigests := []string{"untagged", "foo"}

		actual := getImageList(createFakeImages(), false, []string{"bar"}, []*regexp.Regexp{})

		err := assertImages(expectedDigests, actual)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Exclude tag pattern", func(t *testing.T) {
		expectedDigests := []string{"untagged"}

		re := regexp.MustCompile("^foo")

		actual := getImageList(createFakeImages(), false, []string{}, []*regexp.Regexp{re})

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
