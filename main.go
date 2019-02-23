package main

import (
	"flag"
	"fmt"

	"github.com/graillus/gcrgc/gcloud"
)

// Settings contains app Configuration
type Settings struct {
	Repository string
	Images     []string
	Date       string
}

// ParseArgs parses compmand-line args and returns a Settings instance
func ParseArgs() *Settings {
	settings := Settings{}
	const (
		defaultRepository = ""
	)

	flag.StringVar(&settings.Repository, "repository", "", "Name of the target repository")
	flag.StringVar(&settings.Repository, "image", "", "Name of the taget image. Mutliple values allowed")

	flag.Parse()

	return &settings
}

func main() {
	settings := ParseArgs()

	fmt.Println(settings)

	repoName := "eu.gcr.io/devfactory-etsglobal"
	repo := gcloud.NewRepository(repoName)

	fmt.Printf("Fetching images in repository [%s]\n", repoName)
	images := repo.ListImages()
	fmt.Printf("%d images found\n", len(images))
	for i := 0; i < len(images); i++ {
		fmt.Println(images[i].Name)
	}

	image := gcloud.NewImage("eu.gcr.io/devfactory-etsglobal/nginx")
	tags := image.ListTags()
	for i := 0; i < len(tags); i++ {
		fmt.Println(tags[i].Digest)
	}

}
