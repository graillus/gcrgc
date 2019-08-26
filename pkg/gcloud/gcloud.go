package gcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/graillus/gcrgc/pkg/cmd"
)

// GCloud abstracts the gcloud cli
type GCloud struct {
	e cmd.Executor
}

// NewGCloud creates a new instance of GCloud
func NewGCloud(e cmd.Executor) *GCloud {
	return &GCloud{e}
}

// ListRepositories gets the list of repositories for current registry
func (g GCloud) ListRepositories(registry string) []Repository {
	args := []string{
		"container",
		"images",
		"list",
		strings.Join([]string{"--repository", registry}, "="),
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	cmd := cmd.NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var repos []Repository
	json.Unmarshal(cmd.Stdout.Bytes(), &repos)

	return repos
}

// ListImages gets the list of images for the given repository name
func (g GCloud) ListImages(repo string, minDate string) []Image {
	args := []string{
		"container",
		"images",
		"list-tags",
		repo,
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--sort-by", "TIMESTAMP"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	if minDate != "" {
		args = append(args, strings.Join([]string{"--filter", "timestamp.datetime<'" + minDate + "'"}, "="))
	}

	cmd := cmd.NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var imgs []Image
	json.Unmarshal(cmd.Stdout.Bytes(), &imgs)

	return imgs
}

// DeleteImage deletes an image
func (g GCloud) DeleteImage(repo string, i *Image, dryRun bool) {
	if dryRun {
		i.IsRemoved = true

		return
	}

	args := []string{
		"container",
		"images",
		"delete",
		strings.Join([]string{repo, i.Digest}, "@"),
		"--force-delete-tags",
		"--quiet",
	}

	cmd := cmd.NewCmd("gcloud", args)
	err := g.e.Exec(cmd)
	if err != nil {
		fmt.Printf("Unable to delete image %s: %s\n", i.Digest, err)

		return
	}

	i.IsRemoved = true
}
