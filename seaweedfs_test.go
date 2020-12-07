package fs

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func setup(manager *SeaweedFSManager) {
	_ = manager.DeleteDir("/test")
}

func cleanup(manager *SeaweedFSManager) {
	// cleanup
	_ = manager.DeleteDir("/test")
}

func TestNewSeaweedFSManager(t *testing.T) {
	_, err := NewSeaweedFSManager()
	require.Nil(t, err)
}

func TestSeaweedFSManager_ListDir(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadDir("./test/data/nested", "/test/data/nested")
	require.Nil(t, err)

	valid := false
	files, err := manager.ListDir("/test/data")
	require.Nil(t, err)
	for _, f1 := range files {
		if f1.Name == "nested" && f1.Children != nil {
			for _, f2 := range f1.Children {
				if f2.Name == "nested_test_data.txt" {
					valid = true
				}
			}
		}
	}
	require.True(t, valid)

	cleanup(manager)
}

func TestSeaweedFSManager_UploadFile(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadFile("./test/data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	files, err := manager.ListDir("/test/data")
	require.Nil(t, err)
	valid := false
	for _, file := range files {
		if file.Name == "test_data.txt" {
			valid = true
		}
	}
	require.True(t, valid)

	cleanup(manager)
}

func TestSeaweedFSManager_UploadDir(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadDir("./test/data/nested", "/test/data/nested")
	require.Nil(t, err)

	valid := false
	files, err := manager.ListDir("/test/data")
	require.Nil(t, err)
	for _, f1 := range files {
		if f1.Name == "nested" && f1.Children != nil {
			for _, f2 := range f1.Children {
				if f2.Name == "nested_test_data.txt" {
					valid = true
				}
			}
		}
	}
	require.True(t, valid)

	cleanup(manager)
}

func TestSeaweedFSManager_DownloadFile(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadFile("./test/data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	err = manager.DownloadFile("/test/data/test_data.txt", "./tmp/test_data.txt")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./tmp/test_data.txt")
	require.Nil(t, err)
	require.NotEmpty(t, data)

	cleanup(manager)
}

func TestSeaweedFSManager_DownloadDir(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadDir("./test/data/nested", "/test/data/nested")
	require.Nil(t, err)

	err = manager.DownloadDir("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./test/data/nested/nested_test_data.txt")
	require.Nil(t, err)
	require.NotEmpty(t, data)

	cleanup(manager)
}

func TestSeaweedFSManager_DeleteFile(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadFile("./test/data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	err = manager.DeleteFile("/test/data/test_data.txt")
	require.Nil(t, err)

	files, err := manager.ListDir("/test/data")
	require.Nil(t, err)
	require.Equal(t, len(files), 0)

	cleanup(manager)
}

func TestSeaweedFSManager_DeleteDir(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadDir("./test/data", "/test/data")
	require.Nil(t, err)

	err = manager.DeleteDir("/test/data/nested")
	require.Nil(t, err)

	files, err := manager.ListDir("/test/data")
	require.Nil(t, err)
	valid := true
	for _, file := range files {
		if file.Name == "nested" && file.IsDir {
			valid = false
		}
	}
	require.True(t, valid)

	cleanup(manager)
}
