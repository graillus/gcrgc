package gcloud

import (
	"bytes"
	"os/exec"
)

// Exec execute a gcloud command
// args are the arguments to be passed after the "gcloud" command
func Exec(args []string) ([]byte, error) {
	cmd := exec.Command("gcloud", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.Bytes(), err
}
