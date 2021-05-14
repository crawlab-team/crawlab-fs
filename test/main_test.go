package test

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	var cmd *exec.Cmd
	// before test
	cmd = exec.Command("docker-compose", "up", "-d")
	runCmd(cmd)

	// test
	m.Run()

	// after test
	cmd = exec.Command("docker-compose", "down")
	runCmd(cmd)
}

func runCmd(cmd *exec.Cmd) {
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
