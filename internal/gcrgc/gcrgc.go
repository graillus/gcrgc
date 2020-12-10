package gcrgc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/graillus/gcrgc/internal/docker"
)

// Settings contains app Configuration
type Settings struct {
	Registry             string
	Repositories         []string
	Date                 time.Time
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
	repositories := gcr.ListRepositories(app.settings.Registry)

	registry := docker.NewRegistry(app.settings.Registry, repositories)
	fmt.Printf("%d repositories found\n", len(repositories))

	tasks := getTaskList(
		gcr,
		getRepositoryList(registry, app.settings),
		app.settings,
	)

	doDelete(gcr, tasks, app.settings.DryRun)
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

// Browse the registry to find out which images have to be deleted
func getTaskList(gcr docker.Provider, repository []docker.Repository, s *Settings) taskList {
	tasks := make(taskList)

	// By default all images pushed before the current time will be taken into account
	var date = time.Now()
	if s.Date != (time.Time{}) {
		// If date setting is set we use its value instead
		date = s.Date
	}

	// Create filters to decide which images should be deleted
	filters := []ImageFilter{
		UntaggedFilter{s.UntaggedOnly},
		TagNameFilter{s.ExcludedTags},
		NewTagNameRegexFilter(s.ExcludedTagPatterns),
		NewSemVerTagNameFilter(s.ExcludeSemVerTags),
	}

	for _, repo := range repository {
		imgs := gcr.ListImages(repo.Name, date)

		filteredImgs := filterImages(imgs, filters)

		tasks[repo.Name] = filteredImgs
		fmt.Printf("%d matches for repository [%s]\n", len(filteredImgs), repo.Name)
	}

	return tasks
}

// Delete the images from remote registry
func doDelete(gcr docker.Provider, tasks taskList, dryRun bool) {
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
			gcr.DeleteImage(k, &i, dryRun)

			r.reportImage(i)
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted images: %d/%d\n", r.TotalDeleted(), r.Total())
}
