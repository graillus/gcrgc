package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/graillus/gcrgc/gcloud"
)

func main() {
	settings := ParseArgs()

	repo := gcloud.NewRepository(settings.Repository)

	fmt.Printf("Fetching images in repository [%s]\n", settings.Repository)
	images := repo.ListImages()
	fmt.Printf("%d images found\n", len(images))

	var includedImages []gcloud.Image
	var notFoundImages []string
	if settings.AllImages == true {
		includedImages, notFoundImages = excludeImages(images, settings)
	} else {
		includedImages, notFoundImages = includeImages(images, settings)
	}

	if len(notFoundImages) > 0 {
		os.Exit(1)
	}

	var total, totalDeleted = 0, 0
	for _, image := range includedImages {
		var filteredTags []gcloud.Tag
		tags := image.ListTags(settings.Date)
		for _, tag := range tags {
			if settings.UntaggedOnly && tag.IsTagged() {
				continue
			}
			exclude := false
			if len(settings.ExcludedTags) > 0 {
				for _, excl := range settings.ExcludedTags {
					if tag.ContainsTag(excl) {
						exclude = true
					}
				}
			}
			if !exclude {
				filteredTags = append(filteredTags, tag)
			}
		}

		fmt.Printf("%d matches for image [%s]\n", len(filteredTags), image.Name)
		for _, tag := range filteredTags {
			total++
			fmt.Printf("Deleting %s %s\n", tag.Digest, strings.Join(tag.Tags, ", "))
			tag.Delete(image.Name, settings.DryRun)
			if tag.IsRemoved {
				totalDeleted++
			}
		}
	}

	fmt.Printf("Done\n\n")
	fmt.Printf("Deleted tags: %d/%d\n", totalDeleted, total)
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
