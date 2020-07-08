package gcrgc

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/graillus/gcrgc/pkg/docker"
)

// Settings contains app Configuration
type Settings struct {
	Registry             string
	Repositories         []string
	Date                 *time.Time
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
	auth := createAuthenticator()
	gcr := NewGCR(auth)

	fmt.Printf("Fetching repositories in registry [%s]\n", app.settings.Registry)
	repos := gcr.ListRepositories(app.settings.Registry)
	registry := docker.NewRegistry(app.settings.Registry, repos)
	fmt.Printf("%d repositories found\n", len(repos))

	filteredRepos := getRepoList(registry, app.settings)
	tasks := getTaskList(gcr, filteredRepos, app.settings)

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
			gcr.DeleteImage(k, &i, app.settings.DryRun)
			r.reportImage(i)
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted images: %d/%d\n", r.TotalDeleted(), r.Total())
}

func createAuthenticator() authn.Authenticator {
	var auth authn.Authenticator
	var err error
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		auth, err = google.NewEnvAuthenticator()
		if err != nil {
			log.Fatalf("Cannot create authenticator: %s\n", err)
		}

		return auth
	}
	auth, err = google.NewGcloudAuthenticator()
	if err != nil {
		log.Fatalf("Cannot create authenticator: %s\n", err)
	}

	return auth
}

type taskList map[string][]docker.Image

func getTaskList(gcr docker.Provider, repos []docker.Repository, s *Settings) taskList {
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

	// By default all images pushed before the current time will be taken into account
	var date = time.Now()
	if s.Date != nil {
		// If date setting is set we use its value instead
		date = *s.Date
	}

	for _, repo := range repos {
		imgs := gcr.ListImages(repo.Name, date)

		filteredImgs := getImageList(imgs, s.UntaggedOnly, s.ExcludedTags, exclTagRegexps)

		tasks[repo.Name] = filteredImgs
		fmt.Printf("%d matches for repository [%s]\n", len(filteredImgs), repo.Name)
	}

	return tasks
}
