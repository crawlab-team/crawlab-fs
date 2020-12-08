package fs

import (
	"github.com/crawlab-team/crawlab-fs/lib/copy"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
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

func TestSeaweedFSManager_GetFile(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadFile("./test/data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	data, err := manager.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is a test data", string(data))

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
	require.Equal(t, 0, len(files))

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

func TestSeaweedFSManager_SyncLocalToRemote(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = copy.CopyDirectory("./test/data", "./tmp/data")
	require.Nil(t, err)

	err = manager.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	data, err := manager.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is a test data", string(data))

	err = ioutil.WriteFile("./tmp/data/test_data.txt", []byte("this is changed data"), os.ModePerm)
	require.Nil(t, err)

	err = manager.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	data, err = manager.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is changed data", string(data))

	err = os.Remove("./tmp/data/test_data.txt")
	require.Nil(t, err)

	err = manager.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	valid := true
	files, err := manager.ListDir("/test/data")
	for _, file := range files {
		if file.Name == "test_data.txt" {
			valid = false
		}
	}
	require.True(t, valid)

	cleanup(manager)
}

func TestSeaweedFSManager_SyncRemoteToLocal(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	if _, err := os.Stat("./tmp/data"); err == nil {
		err = os.RemoveAll("./tmp/data")
		require.Nil(t, err)
	}

	err = manager.UploadDir("./test/data", "/test/data")
	require.Nil(t, err)

	err = manager.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./tmp/data/test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is a test data", string(data))

	err = manager.UpdateFile("/test/data/test_data.txt", []byte("this is changed data"))
	require.Nil(t, err)

	err = manager.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err = ioutil.ReadFile("./tmp/data/test_data.txt")
	require.Equal(t, "this is changed data", string(data))

	err = manager.DeleteFile("/test/data/test_data.txt")
	require.Nil(t, err)

	err = manager.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	_, err = os.Stat("./tmp/data/test_data.txt")
	require.NotNil(t, err)

	cleanup(manager)
}

func TestSeaweedFSManager_UpdateFile(t *testing.T) {
	manager, err := NewSeaweedFSManager()
	require.Nil(t, err)

	setup(manager)

	err = manager.UploadDir("./test/data", "/test/data")
	require.Nil(t, err)

	err = manager.UpdateFile("/test/data/test_data.txt", []byte("this is changed data"))
	require.Nil(t, err)

	data, err := manager.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is changed data", string(data))

	cleanup(manager)
}
