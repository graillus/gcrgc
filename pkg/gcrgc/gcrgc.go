package gcrgc

import (
	"fmt"
	"os"
	"strings"

	"github.com/graillus/gcrgc/pkg/cmd"
	"github.com/graillus/gcrgc/pkg/gcloud"
)

// Settings contains app Configuration
type Settings struct {
	Registry             string
	Repositories         stringList
	Date                 string
	UntaggedOnly         bool
	DryRun               bool
	AllRepositories      bool
	ExcludedRepositories stringList
	ExcludedTags         stringList
}

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)

	return nil
}

// App is the main application
type App struct {
	settings *Settings
}

// NewApp creates a new instance of the application
func NewApp(s *Settings) *App {
	return &App{s}
}

// Start the application
func (app *App) Start() {
	cli := cmd.NewCli()
	gcloudCmd := gcloud.NewGCloud(cli)

	registry := gcloud.NewRegistry(app.settings.Registry)

	fmt.Printf("Fetching repositories in registry [%s]\n", registry.Name)
	repos := gcloudCmd.ListRepositories(registry.Name)
	fmt.Printf("%d repositories found\n", len(repos))

	registry.Repositories = repos

	r := newReport()
	tasks := getTaskList(gcloudCmd, *registry, app.settings)
	for k, v := range tasks {
		fmt.Printf("Cleaning repository [%s] (%d matches)\n", k, len(v))
		for _, i := range v {
			fmt.Printf("Deleting %s %s\n", i.Digest, strings.Join(i.Tags, ", "))
			gcloudCmd.DeleteImage(k, &i, app.settings.DryRun)
			r.reportImage(i)
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted images: %d/%d\n", r.TotalDeleted(), r.Total())
}

type images []gcloud.Image
type taskList map[string]images

func getTaskList(gcloudCmd *gcloud.GCloud, registry gcloud.Registry, s *Settings) taskList {
	var included []gcloud.Repository
	var notFound []string
	if s.AllRepositories == true {
		included, notFound = excludeRepos(registry, s)
	} else {
		included, notFound = includeRepos(registry, s)
	}

	if len(notFound) > 0 {
		os.Exit(1)
	}

	tasks := make(taskList)
	for _, repo := range included {
		var filteredImgs images
		imgs := gcloudCmd.ListImages(repo.Name, s.Date)
		for _, img := range imgs {
			if s.UntaggedOnly && img.IsTagged() {
				continue
			}
			exclude := false
			if len(s.ExcludedTags) > 0 {
				for _, excl := range s.ExcludedTags {
					if img.ContainsTag(excl) {
						exclude = true
					}
				}
			}
			if !exclude {
				filteredImgs = append(filteredImgs, img)
			}
		}
		tasks[repo.Name] = filteredImgs
		fmt.Printf("%d matches for repository [%s]\n", len(filteredImgs), repo.Name)
	}

	return tasks
}

func excludeRepos(registry gcloud.Registry, settings *Settings) ([]gcloud.Repository, []string) {
	var notFound []string
	included := registry.Repositories
	for _, img := range settings.ExcludedRepositories {
		repoName := settings.Registry + "/" + img
		if !registry.ContainsRepository(repoName) {
			notFound = append(notFound, repoName)
			fmt.Printf("Cannot exclude repository [%s]: it does not exist in this registry\n", repoName)
		} else {
			for i := 0; i < len(included); i++ {
				if included[i].Name == repoName {
					included = append(included[:i], included[i+1:]...)
				}
			}
		}
	}

	return included, notFound
}

func includeRepos(registry gcloud.Registry, settings *Settings) ([]gcloud.Repository, []string) {
	var notFound []string
	var included []gcloud.Repository
	for _, img := range settings.Repositories {
		repoName := settings.Registry + "/" + img
		if !registry.ContainsRepository(repoName) {
			notFound = append(notFound, repoName)
			fmt.Printf("Cannot include repository [%s]: it does not exist in this registry\n", repoName)
		} else {
			included = append(included, *gcloud.NewRepository(repoName))
		}
	}

	return included, notFound
}
