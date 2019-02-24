package gcloud

import (
	"encoding/json"
	"log"
	"strings"
)

// Repository wraps the gcloud container images command
type Repository struct {
	name string
}

// NewRepository returns a new instance of Repository
func NewRepository(name string) *Repository {
	return &Repository{name}
}

// ListImages gets the list of repositories for current registry
func (r Repository) ListImages() []Image {
	parts := []string{
		"container",
		"images",
		"list",
		strings.Join([]string{"--repository", r.name}, "="),
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	output, err := Exec(parts)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var images []Image
	json.Unmarshal(output, &images)

	return images
}
