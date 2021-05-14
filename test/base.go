package test

import (
	fs "github.com/crawlab-team/crawlab-fs"
	"testing"
	"time"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	m fs.Manager
}

func (t *Test) Setup(t2 *testing.T) {
	t.Cleanup()
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = T.m.DeleteDir("/test")

	// wait to avoid caching
	time.Sleep(200 * time.Millisecond)
}

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// manager
	t.m, err = fs.NewSeaweedFsManager(
		fs.WithFilerUrl("http://localhost:8888"),
		fs.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}
