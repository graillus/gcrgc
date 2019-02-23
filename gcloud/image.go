package gcloud

import (
	"encoding/json"
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
func (i Image) ListTags() []Tag {
	parts := []string{
		"container",
		"images",
		"list-tags",
		i.Name,
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--filter", "\"timestamp.datetime < '2019-01-01'\""}, "="),
		strings.Join([]string{"--sort-by", "TIMESTAMP"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	output := Exec(parts)

	var tags []Tag
	json.Unmarshal(output, &tags)

	return tags
}
