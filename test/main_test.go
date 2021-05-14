package test

import (
	"testing"
)

func TestMain(m *testing.M) {
	// before test
	filePath, err := StartTestContainers()
	if err != nil {
		panic(err)
	}

	// test
	m.Run()

	// after test
	_ = StopTestContainers(filePath)
}
