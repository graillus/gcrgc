package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// GCloud abstracts the gcloud cli
type GCloud struct {
	e CmdExecutor
}

// NewGCloud creates a new instance of GCloud
func NewGCloud(e CmdExecutor) *GCloud {
	return &GCloud{e}
}

// ListImages gets the list of repositories for current registry
func (g GCloud) ListImages(repository string) []Image {
	args := []string{
		"container",
		"images",
		"list",
		strings.Join([]string{"--repository", repository}, "="),
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	cmd := NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var images []Image
	json.Unmarshal(cmd.Stdout.Bytes(), &images)

	return images
}

// ListTags gets the list of tags for the current image
func (g GCloud) ListTags(image string, minDate string) []Tag {
	args := []string{
		"container",
		"images",
		"list-tags",
		image,
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--sort-by", "TIMESTAMP"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	if minDate != "" {
		args = append(args, strings.Join([]string{"--filter", "timestamp.datetime<'" + minDate + "'"}, "="))
	}

	cmd := NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var tags []Tag
	json.Unmarshal(cmd.Stdout.Bytes(), &tags)

	return tags
}

// Delete deletes the tag
func (g GCloud) Delete(image string, t *Tag, dryRun bool) {
	if dryRun {
		t.IsRemoved = true

		return
	}

	args := []string{
		"container",
		"images",
		"delete",
		strings.Join([]string{image, t.Digest}, "@"),
		"--force-delete-tags",
		"--quiet",
	}

	cmd := NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		fmt.Printf("Unable to delete tag %s: %s\n", t.Digest, err)

		return
	}

	t.IsRemoved = true
}
