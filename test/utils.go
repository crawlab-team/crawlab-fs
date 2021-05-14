package test

import (
	"fmt"
	"github.com/crawlab-team/go-trace"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func StartTestContainers() (filePath string, err error) {
	// docker-compose.yml
	tmpDir := os.TempDir()
	id, _ := uuid.NewUUID()
	fileName := fmt.Sprintf("crawlab_fs_docker-compose_%s.yml", id)
	data, err := Asset("docker-compose.yml")
	if err != nil {
		return filePath, trace.TraceError(err)
	}
	filePath = path.Join(tmpDir, fileName)
	if err := ioutil.WriteFile(filePath, data, os.FileMode(0766)); err != nil {
		return filePath, trace.TraceError(err)
	}

	// docker-compose up
	cmd := exec.Command("docker-compose", "up", "-d", "-f", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return filePath, nil
}

func StopTestContainers(filePath string) (err error) {
	// docker-compose down
	cmd := exec.Command("docker-compose", "down", "-f", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		return trace.TraceError(err)
	}

	return nil
}
