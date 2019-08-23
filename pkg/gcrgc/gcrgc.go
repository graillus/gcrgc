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
	Repository     string
	Images         stringList
	Date           string
	UntaggedOnly   bool
	DryRun         bool
	AllImages      bool
	ExcludedImages stringList
	ExcludedTags   stringList
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
	gcloud := gcloud.NewGCloud(cli)

	fmt.Printf("Fetching images in repository [%s]\n", app.settings.Repository)
	images := gcloud.ListImages(app.settings.Repository)
	fmt.Printf("%d images found\n", len(images))

	r := newReport()
	tasks := getTaskList(gcloud, images, app.settings)
	for k, v := range tasks {
		fmt.Printf("Deleting %d matches for image [%s]\n", len(v), k)
		for _, t := range v {
			fmt.Printf("Deleting %s %s\n", t.Digest, strings.Join(t.Tags, ", "))
			gcloud.Delete(k, &t, app.settings.DryRun)
			r.reportTag(t)
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted tags: %d/%d\n", r.TotalDeleted(), r.Total())
}

type tags []gcloud.Tag
type taskList map[string]tags

func getTaskList(gCloud *gcloud.GCloud, images []gcloud.Image, s *Settings) taskList {
	var includedImages []gcloud.Image
	var notFoundImages []string
	if s.AllImages == true {
		includedImages, notFoundImages = excludeImages(images, s)
	} else {
		includedImages, notFoundImages = includeImages(images, s)
	}

	if len(notFoundImages) > 0 {
		os.Exit(1)
	}

	tasks := make(taskList)
	for _, image := range includedImages {
		var filteredTags tags
		tags := gCloud.ListTags(image.Name, s.Date)
		for _, tag := range tags {
			if s.UntaggedOnly && tag.IsTagged() {
				continue
			}
			exclude := false
			if len(s.ExcludedTags) > 0 {
				for _, excl := range s.ExcludedTags {
					if tag.ContainsTag(excl) {
						exclude = true
					}
				}
			}
			if !exclude {
				filteredTags = append(filteredTags, tag)
			}
		}
		tasks[image.Name] = filteredTags
		fmt.Printf("%d matches for image [%s]\n", len(filteredTags), image.Name)
	}

	return tasks
}

func excludeImages(images []gcloud.Image, settings *Settings) ([]gcloud.Image, []string) {
	var notFoundImages []string
	includedImages := images
	for _, img := range settings.ExcludedImages {
		imageName := settings.Repository + "/" + img
		if !gcloud.ContainsImage(imageName, images) {
			notFoundImages = append(notFoundImages, imageName)
			fmt.Printf("Cannot exclude image [%s]: it does not exist in this repository\n", imageName)
		} else {
			for i := 0; i < len(includedImages); i++ {
				if includedImages[i].Name == imageName {
					includedImages = append(includedImages[:i], includedImages[i+1:]...)
				}
			}
		}
	}

	return includedImages, notFoundImages
}

func includeImages(images []gcloud.Image, settings *Settings) ([]gcloud.Image, []string) {
	var includedImages []gcloud.Image
	var notFoundImages []string
	for _, img := range settings.Images {
		imageName := settings.Repository + "/" + img
		if !gcloud.ContainsImage(imageName, images) {
			notFoundImages = append(notFoundImages, imageName)
			fmt.Printf("Cannot include image [%s]: it does not exist in this repository\n", imageName)
		} else {
			includedImages = append(includedImages, *gcloud.NewImage(imageName))
		}
	}

	return includedImages, notFoundImages
}
