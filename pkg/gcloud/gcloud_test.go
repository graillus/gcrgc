package gcloud

import (
	"strings"
	"testing"

	"github.com/graillus/gcrgc/pkg/cmd"
)

type cliMock struct {
	cmd *cmd.Cmd
	ret error
}

func (c *cliMock) Exec(gcmd *cmd.Cmd) error {
	c.cmd = gcmd

	return c.ret
}

func TestListRepositories(t *testing.T) {
	cli := cliMock{nil, nil}
	gcloud := NewGCloud(&cli)

	gcloud.ListRepositories("gcr.io/project-id")

	if cli.cmd.Name != "gcloud" {
		t.Errorf("Expected command name to be %s, but got %s", "gcloud", cli.cmd.Name)
	}

	expectedCommand := "container images list --repository=gcr.io/project-id --format=json --limit=999999"
	actualCommand := strings.Join(cli.cmd.Args, " ")
	if actualCommand != expectedCommand {
		t.Errorf("Expected command to be \"%s\", but got \"%s\"", expectedCommand, actualCommand)
	}
}

func TestListImages(t *testing.T) {
	cli := cliMock{nil, nil}
	gcloud := NewGCloud(&cli)

	gcloud.ListImages("gcr.io/project-id/image", "2019-01-01")

	if cli.cmd.Name != "gcloud" {
		t.Errorf("Expected command name to be %s, but got %s", "gcloud", cli.cmd.Name)
	}

	expectedCommand := "container images list-tags gcr.io/project-id/image --format=json --sort-by=TIMESTAMP --limit=999999 --filter=timestamp.datetime<'2019-01-01'"
	actualCommand := strings.Join(cli.cmd.Args, " ")
	if actualCommand != expectedCommand {
		t.Errorf("Expected command to be \"%s\", but got \"%s\"", expectedCommand, actualCommand)
	}
}

func TestDeleteImage(t *testing.T) {
	cli := cliMock{nil, nil}
	gcloud := NewGCloud(&cli)
	img := Image{Digest: "sha256:digest"}

	gcloud.DeleteImage("gcr.io/project-id/image", &img, false)

	if cli.cmd.Name != "gcloud" {
		t.Errorf("Expected command name to be %s, but got %s", "gcloud", cli.cmd.Name)
	}

	expectedCommand := "container images delete gcr.io/project-id/image@sha256:digest --force-delete-tags --quiet"
	actualCommand := strings.Join(cli.cmd.Args, " ")
	if actualCommand != expectedCommand {
		t.Errorf("Expected command to be \"%s\", but got \"%s\"", expectedCommand, actualCommand)
	}

	if img.IsRemoved != true {
		t.Errorf("Expected image to be marked as deleted")
	}
}

func TestDeleteImageDryRun(t *testing.T) {
	cli := cliMock{nil, nil}
	gcloud := NewGCloud(&cli)
	img := Image{Digest: "sha256:digest"}

	gcloud.DeleteImage("gcr.io/project-id/image", &img, true)

	if cli.cmd != nil {
		t.Errorf("Unexpected call to gcloud cli")
	}

	if img.IsRemoved != true {
		t.Errorf("Expected image to be marked as deleted")
	}
}
