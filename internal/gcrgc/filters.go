package gcrgc

import (
	"fmt"
	"regexp"

	"github.com/graillus/gcrgc/internal/docker"
)

func getRepositoryList(registry *docker.Registry, s *Settings) []docker.Repository {
	var repos []docker.Repository

	// Make a list of repositories to exclude prefixed with registry name
	exclRepos := []string{}
	for _, r := range s.ExcludedRepositories {
		exclRepos = append(exclRepos, repositoryName(registry.Name, r))
	}

	// Make a list of repositories to include prefixed with registry name
	inclRepos := []string{}
	for _, r := range s.Repositories {
		inclRepos = append(inclRepos, repositoryName(registry.Name, r))
	}

	if s.AllRepositories == true {
		repos = excludeRepositories(registry, exclRepos)
	} else {
		repos = includeRepositories(registry, inclRepos)
	}

	return repos
}

func repositoryName(registry string, image string) string {
	return registry + "/" + image
}

func excludeRepositories(registry *docker.Registry, toExclude []string) []docker.Repository {
	included := registry.Repositories

	for _, repoName := range toExclude {
		if !registry.ContainsRepository(repoName) {
			fmt.Printf("Warning: Cannot exclude repository [%s]: it does not exist in this registry\n", repoName)

			continue
		}

		for i := 0; i < len(included); i++ {
			if included[i].Name == repoName {
				included = append(included[:i], included[i+1:]...)
			}
		}
	}

	return included
}

func includeRepositories(registry *docker.Registry, toInclude []string) []docker.Repository {
	var included []docker.Repository

	for _, repoName := range toInclude {
		if !registry.ContainsRepository(repoName) {
			fmt.Printf("Warning: Cannot include repository [%s]: it does not exist in this registry\n", repoName)

			continue
		}

		included = append(included, *docker.NewRepository(repoName))
	}

	return included
}

func filterImages(imgs []docker.Image, filters []ImageFilter) []docker.Image {
	if len(filters) == 0 {
		return imgs
	}

	var list []docker.Image
	for _, img := range imgs {
		// All filters must return true for the image to be eligible for deletion
		eligible := true
		for _, f := range filters {
			eligible = f.Apply(&img) && eligible
		}

		if eligible {
			list = append(list, img)
		}
	}

	return list
}

// ImageFilter is an interface for docker image filters
type ImageFilter interface {
	// Apply should return true if the image should be planned for deletion
	Apply(i *docker.Image) bool
}

// UntaggedFilter filters images that have no tag
type UntaggedFilter struct {
	enabled bool
}

// Apply returns true when the image has no tags
func (f UntaggedFilter) Apply(img *docker.Image) bool {
	if f.enabled == false {
		// Filter is disabled, image is eligible for deletion
		return true
	}

	return !img.IsTagged()
}

// TagNameFilter compares the image tags against the exclusion list.
type TagNameFilter struct {
	excluded []string
}

// Apply returns true if if no tag was in the exclusion list.
func (f TagNameFilter) Apply(img *docker.Image) bool {
	for _, tag := range f.excluded {
		if img.HasTag(tag) {
			return false
		}
	}

	return true
}

// TagNameRegexFilter compares the image tags against regular expressions from the exclusion list.
type TagNameRegexFilter struct {
	excluded []*regexp.Regexp
}

// NewTagNameRegexFilter creates a new SemVerTagNameFilter
func NewTagNameRegexFilter(patterns []string) *TagNameRegexFilter {
	var re []*regexp.Regexp

	for _, p := range patterns {
		re = append(re, regexp.MustCompile(p))
	}

	return &TagNameRegexFilter{re}
}

// Apply returns true if if no tag was in the exclusion list.
func (f TagNameRegexFilter) Apply(img *docker.Image) bool {
	for _, re := range f.excluded {
		if img.HasTagRegexp(re) {
			return false
		}
	}

	return true
}

// SemVerTagNameFilter compares the image tags against regular expressions from the exclusion list.
type SemVerTagNameFilter struct {
	enabled bool
	regex   *regexp.Regexp
}

// NewSemVerTagNameFilter creates a new SemVerTagNameFilter
func NewSemVerTagNameFilter(enabled bool) *SemVerTagNameFilter {
	const semverPattern = `^[vV]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

	return &SemVerTagNameFilter{
		enabled,
		regexp.MustCompile(semverPattern),
	}
}

// Apply returns true if if no tag was in the exclusion list.
func (f *SemVerTagNameFilter) Apply(img *docker.Image) bool {
	if !f.enabled {
		// Filter is disabled, image is eligible for deletion.
		return true
	}

	return !img.HasTagRegexp(f.regex)
}
