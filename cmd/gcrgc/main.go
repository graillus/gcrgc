package main

import (
	"fmt"
	"os"
	"strings"
)

type tags []Tag
type taskList map[string]tags

func main() {
	settings := ParseArgs()

	cli := NewCli()
	gcloud := NewGCloud(cli)

	fmt.Printf("Fetching images in repository [%s]\n", settings.Repository)
	images := gcloud.ListImages(settings.Repository)
	fmt.Printf("%d images found\n", len(images))

	r := newReport()
	tasks := getTaskList(gcloud, images, settings)
	for k, v := range tasks {
		fmt.Printf("Deleting %d matches for image [%s]\n", len(v), k)
		for _, t := range v {
			fmt.Printf("Deleting %s %s\n", t.Digest, strings.Join(t.Tags, ", "))
			gcloud.Delete(k, &t, settings.DryRun)
			r.reportTag(t)
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted tags: %d/%d\n", r.TotalDeleted(), r.Total())
}

func getTaskList(gcloud *GCloud, images []Image, s *Settings) taskList {
	var includedImages []Image
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
		tags := gcloud.ListTags(image.Name, s.Date)
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

func excludeImages(images []Image, settings *Settings) ([]Image, []string) {
	var notFoundImages []string
	includedImages := images
	for _, img := range settings.ExcludedImages {
		imageName := settings.Repository + "/" + img
		if !ContainsImage(imageName, images) {
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

func includeImages(images []Image, settings *Settings) ([]Image, []string) {
	var includedImages []Image
	var notFoundImages []string
	for _, img := range settings.Images {
		imageName := settings.Repository + "/" + img
		if !ContainsImage(imageName, images) {
			notFoundImages = append(notFoundImages, imageName)
			fmt.Printf("Cannot include image [%s]: it does not exist in this repository\n", imageName)
		} else {
			includedImages = append(includedImages, *NewImage(imageName))
		}
	}

	return includedImages, notFoundImages
}
