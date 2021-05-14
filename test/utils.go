package test

import (
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/go-trace"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

func StartTestContainers() (filePath string, err error) {
	// force remove containers
	removeTestContainers()

	// directories
	tmpDir := os.TempDir()
	id, _ := uuid.NewUUID()
	dirName := id.String()
	dirPath := path.Join(tmpDir, dirName)
	if _, err := os.Stat(dirPath); err != nil {
		if err := os.MkdirAll(dirPath, os.FileMode(0766)); err != nil {
			return filePath, trace.TraceError(err)
		}
	}

	// write to docker-compose.yml
	data, err := Asset("docker-compose.yml")
	if err != nil {
		return filePath, trace.TraceError(err)
	}
	filePath = path.Join(dirPath, "docker-compose.yml")
	if err := ioutil.WriteFile(filePath, data, os.FileMode(0766)); err != nil {
		return filePath, trace.TraceError(err)
	}

	// docker-compose up
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// wait for containers to be ready
	time.Sleep(5 * time.Second)
	err = backoff.RetryNotify(func() error {
		_, err := T.m.ListDir("/", true)
		if err != nil {
			return err
		}
		return nil
	}, backoff.NewConstantBackOff(5*time.Second), func(err error, duration time.Duration) {
		log.Infof("seaweedfs containers not ready, re-attempt in %.1f seconds", duration.Seconds())
	})
	if err != nil {
		return filePath, trace.TraceError(err)
	}

	return filePath, nil
}

func StopTestContainers(filePath string) (err error) {
	// directories
	dirPath := path.Dir(filePath)

	// docker-compose down
	cmd := exec.Command("docker-compose", "down", "-v")
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		return trace.TraceError(err)
	}

	// remove networks
	_ = exec.Command("docker", "network", "prune")

	return nil
}

func removeTestContainers() {
	_ = exec.Command("docker", "rm", "-f", "seaweedfs-master").Run()
	_ = exec.Command("docker", "rm", "-f", "seaweedfs-volume").Run()
	_ = exec.Command("docker", "rm", "-f", "seaweedfs-filer").Run()
}
