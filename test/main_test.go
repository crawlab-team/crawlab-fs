package test

import (
	"testing"
)

func TestMain(m *testing.M) {
	// before test
	if err := StartTestSeaweedFs(); err != nil {
		panic(err)
	}

	// test
	m.Run()

	// after test
	if err := StopTestSeaweedFs(); err != nil {
		panic(err)
	}
}
