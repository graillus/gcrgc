package gcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/graillus/gcrgc/pkg/cmd"
	"github.com/graillus/gcrgc/pkg/docker"
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
func (g GCloud) ListRepositories(registry string) []docker.Repository {
	args := []string{
		"container",
		"images",
		"list",
		strings.Join([]string{"--repository", registry}, "="),
		strings.Join([]string{"--format", "json"}, "="),
		strings.Join([]string{"--limit", "999999"}, "="),
	}

	cmd, err := g.exec(args)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var reposData []repository
	json.Unmarshal(cmd.Stdout.Bytes(), &reposData)

	var repos []docker.Repository
	for _, repoData := range reposData {
		repos = append(repos, *docker.NewRepository(repoData.Name))
	}

	return repos
}

// ListImages gets the list of images for the given repository name
func (g GCloud) ListImages(repo string, minDate string) []docker.Image {
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

	cmd, err := g.exec(args)
	if err != nil {
		log.Fatalf("Command failed with %s\n", err)
	}

	var imgsData []image
	json.Unmarshal(cmd.Stdout.Bytes(), &imgsData)

	var imgs []docker.Image
	for _, imgData := range imgsData {
		imgs = append(imgs, *docker.NewImage(imgData.Digest, imgData.Tags))
	}

	return imgs
}

// DeleteImage deletes an image
func (g GCloud) DeleteImage(repo string, i *docker.Image, dryRun bool) {
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

	_, err := g.exec(args)
	if err != nil {
		fmt.Printf("Unable to delete image %s: %s\n", i.Digest, err)

		return
	}

	i.IsRemoved = true
}

func (g GCloud) exec(args []string) (*cmd.Cmd, error) {
	cmd := cmd.NewCmd("gcloud", args)
	err := g.e.Exec(cmd)

	return cmd, err
}

// Repository represents a docker image inside repository
type repository struct {
	Name string `json:"name"`
}

// Image represents a repository image
type image struct {
	Digest    string    `json:"digest"`
	Tags      []string  `json:"tags"`
	Timestamp timestamp `json:"timestamp"`
}

// Timestamp holds the image's date and time information
type timestamp struct {
	Datetime string `json:"datetime"`
}
