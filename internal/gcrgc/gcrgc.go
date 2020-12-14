package gcrgc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
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

	// Report deleted images
	r := newReport()
	for _, task := range tasks {
		r.reportImage(*task.image)
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

type task struct {
	repository string
	image      *docker.Image
}

// Browse the registry to find out which images have to be deleted
func getTaskList(gcr docker.Provider, repository []docker.Repository, s *Settings) []task {
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

	tasks := []task{}

	for _, repo := range repository {
		imgs := gcr.ListImages(repo.Name, date)

		filteredImgs := filterImages(imgs, filters)

		fmt.Printf("%d matches for repository [%s]\n", len(filteredImgs), repo.Name)

		for _, i := range filteredImgs {
			img := i
			tasks = append(tasks, task{repo.Name, &img})
		}
	}

	return tasks
}

// Delete the images from remote registry
func doDelete(gcr docker.Provider, tasks []task, dryRun bool) {
	if len(tasks) == 0 {
		fmt.Println("Nothing to do.")

		return
	}

	// Create a pool of 8 workers. The number is arbitrary, but it seems that inscreasing the worker count doesn't affect
	// the performance since the google container registry API has some kind of per-user rate limiting.
	wp := workerpool.New(8)

	// Let's submit the tasks to the workers pool
	for _, t := range tasks {
		task := t
		wp.Submit(func() {
			fmt.Printf("Deleting %s %s\n", task.image.Digest, strings.Join(task.image.Tags, ", "))
			gcr.DeleteImage(task.repository, task.image, dryRun)
		})
	}
	wp.StopWait()
}
