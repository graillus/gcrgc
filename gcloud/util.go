package gcloud

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// Exec execute a gcloud command
// args are the arguments to be passed after the "gcloud" command
func Exec(args []string) []byte {
	cmd := exec.Command("gcloud", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Println(cmd.Args)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed with %s\n%s", err, string(stderr.Bytes()))
	}

	return stdout.Bytes()
}
