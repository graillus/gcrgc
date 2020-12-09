package gcrgc

import (
	"fmt"
	"regexp"

	"github.com/graillus/gcrgc/internal/docker"
)

func getRepoList(registry *docker.Registry, s *Settings) []docker.Repository {
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
		repos = excludeRepos(registry, exclRepos)
	} else {
		repos = includeRepos(registry, inclRepos)
	}

	return repos
}

func repositoryName(registry string, image string) string {
	return registry + "/" + image
}

func excludeRepos(registry *docker.Registry, toExclude []string) []docker.Repository {
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

func includeRepos(registry *docker.Registry, toInclude []string) []docker.Repository {
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

func getImageList(imgs []docker.Image, untaggedOnly bool, exludedTags []string, exclTagRegexps []*regexp.Regexp) []docker.Image {
	var filteredImgs []docker.Image
	for _, img := range imgs {
		if untaggedOnly && img.IsTagged() {
			continue
		}

		exclude := false
		for _, excl := range exludedTags {
			if img.HasTag(excl) {
				exclude = true
			}
		}

		for _, excl := range exclTagRegexps {
			if img.HasTagRegexp(excl) {
				exclude = true
			}
		}

		if !exclude {
			filteredImgs = append(filteredImgs, img)
		}
	}

	return filteredImgs
}
