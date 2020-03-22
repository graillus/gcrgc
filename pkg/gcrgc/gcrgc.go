package gcrgc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/graillus/gcrgc/pkg/cmd"
	"github.com/graillus/gcrgc/pkg/docker"
	"github.com/graillus/gcrgc/pkg/gcloud"
)

// Settings contains app Configuration
type Settings struct {
	Registry             string
	Repositories         []string
	Date                 string
	UntaggedOnly         bool
	DryRun               bool
	AllRepositories      bool
	ExcludedRepositories []string
	ExcludedTags         []string
	ExcludedTagPatterns  []string
	ExcludeSemVerTags    bool
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

	fmt.Printf("Fetching repositories in registry [%s]\n", app.settings.Registry)
	repos := gcloudCmd.ListRepositories(app.settings.Registry)
	registry := docker.NewRegistry(app.settings.Registry, repos)
	fmt.Printf("%d repositories found\n", len(repos))

	filteredRepos := getRepoList(registry, app.settings)
	tasks := getTaskList(gcloudCmd, filteredRepos, app.settings)

	if len(tasks) == 0 {
		fmt.Println("Nothing to do.")

		return
	}

	r := newReport()
	for k, v := range tasks {
		if len(v) == 0 {
			fmt.Printf("No images to clean in repository [%s]\n", k)

			continue
		}

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

type taskList map[string][]docker.Image

func getTaskList(gcloudCmd docker.Provider, repos []docker.Repository, s *Settings) taskList {
	const semverPattern = `^[vV]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	var exclTagRegexps []*regexp.Regexp

	for _, p := range s.ExcludedTagPatterns {
		regexp := regexp.MustCompile(p)
		exclTagRegexps = append(exclTagRegexps, regexp)
	}

	if s.ExcludeSemVerTags == true {
		regexp := regexp.MustCompile(semverPattern)
		exclTagRegexps = append(exclTagRegexps, regexp)
	}

	tasks := make(taskList)

	for _, repo := range repos {
		imgs := gcloudCmd.ListImages(repo.Name, s.Date)

		filteredImgs := getImageList(imgs, s.UntaggedOnly, s.ExcludedTags, exclTagRegexps)

		tasks[repo.Name] = filteredImgs
		fmt.Printf("%d matches for repository [%s]\n", len(filteredImgs), repo.Name)
	}

	return tasks
}
