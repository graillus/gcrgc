package gcloud

import (
	"encoding/json"
	"log"
	"strings"
)

// Image represents a docker image inside repository
type Image struct {
	Name string `json:"name"`
}

// NewImage returns a new instance of Image
func NewImage(name string) *Image {
	return &Image{name}
}

// ListTags gets the list of tags for the current image
func (i Image) ListTags(minDate string) []Tag {
	parts := []string{
		"container",
		"images",
		"list-tags",
		i.Name,
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--sort-by", "TIMESTAMP"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	if minDate != "" {
		parts = append(parts, strings.Join([]string{"--filter", "timestamp.datetime<'" + minDate + "'"}, "="))
	}

	output, err := Exec(parts)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var tags []Tag
	json.Unmarshal(output, &tags)

	return tags
}

// ContainsImage checks if an image is present in an array of Image structs
func ContainsImage(name string, images []Image) bool {
	for _, item := range images {
		if name == item.Name {
			return true
		}
	}

	return false
}
