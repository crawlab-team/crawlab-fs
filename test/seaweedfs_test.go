package test

import (
	"github.com/crawlab-team/crawlab-fs"
	"github.com/crawlab-team/crawlab-fs/lib/copy"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewSeaweedFsManager(t *testing.T) {
	_, err := fs.NewSeaweedFsManager()
	require.Nil(t, err)
}

func TestSeaweedFsManager_ListDir(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data/nested", "/test/data/nested")
	require.Nil(t, err)

	valid := false
	files, err := T.m.ListDir("/test/data", true)
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
}

func TestSeaweedFsManager_UploadFile(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadFile("./data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	files, err := T.m.ListDir("/test/data", true)
	require.Nil(t, err)
	valid := false
	for _, file := range files {
		if file.Name == "test_data.txt" {
			valid = true
		}
	}
	require.True(t, valid)
}

func TestSeaweedFsManager_UploadDir(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data/nested", "/test/data/nested")
	require.Nil(t, err)

	valid := false
	files, err := T.m.ListDir("/test/data", true)
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
}

func TestSeaweedFsManager_GetFile(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadFile("./data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	data, err := T.m.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is a test data", string(data))
}

func TestSeaweedFsManager_DownloadFile(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadFile("./data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	err = T.m.DownloadFile("/test/data/test_data.txt", "./tmp/test_data.txt")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./tmp/test_data.txt")
	require.Nil(t, err)
	require.NotEmpty(t, data)
}

func TestSeaweedFsManager_DownloadDir(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data/nested", "/test/data/nested")
	require.Nil(t, err)

	err = T.m.DownloadDir("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./data/nested/nested_test_data.txt")
	require.Nil(t, err)
	require.NotEmpty(t, data)
}

func TestSeaweedFsManager_DeleteFile(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadFile("./data/test_data.txt", "/test/data/test_data.txt")
	require.Nil(t, err)

	err = T.m.DeleteFile("/test/data/test_data.txt")
	require.Nil(t, err)

	files, err := T.m.ListDir("/test/data", true)
	require.Nil(t, err)
	require.Equal(t, 0, len(files))
}

func TestSeaweedFsManager_DeleteDir(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data", "/test/data")
	require.Nil(t, err)

	err = T.m.DeleteDir("/test/data/nested")
	require.Nil(t, err)

	files, err := T.m.ListDir("/test/data", true)
	require.Nil(t, err)
	valid := true
	for _, file := range files {
		if file.Name == "nested" && file.IsDir {
			valid = false
		}
	}
	require.True(t, valid)
}

func TestSeaweedFsManager_SyncLocalToRemote(t *testing.T) {
	var err error
	T.Setup(t)

	err = copy.CopyDirectory("./data", "./tmp/data")
	require.Nil(t, err)

	err = T.m.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	data, err := T.m.GetFile("/test/data/test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is a test data", string(data))

	data, err = T.m.GetFile("/test/data/nested/nested_test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is nested test data", string(data))

	err = ioutil.WriteFile("./tmp/data/test_data.txt", []byte("this is changed data"), os.ModePerm)
	require.Nil(t, err)

	err = T.m.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	data, err = T.m.GetFile("/test/data/test_data.txt")
	require.Equal(t, "this is changed data", string(data))

	err = os.Remove("./tmp/data/test_data.txt")
	require.Nil(t, err)

	err = T.m.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)

	valid := true
	files, err := T.m.ListDir("/test/data", true)
	for _, file := range files {
		if file.Name == "test_data.txt" {
			valid = false
		}
	}
	require.True(t, valid)

	// check if directory is deleted after sync
	err = ioutil.WriteFile("./tmp/test.txt", []byte("test"), os.ModePerm)
	require.Nil(t, err)
	err = T.m.UploadFile("./tmp/test.txt", "/test/data/folder1/test.txt")
	require.Nil(t, err)
	time.Sleep(1 * time.Second)
	err = T.m.SyncLocalToRemote("./tmp/data", "/test/data")
	require.Nil(t, err)
	valid = true
	files, err = T.m.ListDir("/test/data", true)
	for _, file := range files {
		if strings.Contains(file.FullPath, "folder1") {
			valid = false
		}
	}
	require.True(t, valid)
}

func TestSeaweedFsManager_SyncRemoteToLocal(t *testing.T) {
	var err error
	T.Setup(t)

	if _, err := os.Stat("./tmp/data"); err == nil {
		err = os.RemoveAll("./tmp/data")
		require.Nil(t, err)
	}

	err = T.m.UploadDir("./data", "/test/data")
	require.Nil(t, err)

	err = T.m.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err := ioutil.ReadFile("./tmp/data/test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is a test data", string(data))

	data, err = ioutil.ReadFile("./tmp/data/nested/nested_test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is nested test data", string(data))

	err = T.m.UpdateFile("/test/data/test_data.txt", []byte("this is changed data"))
	require.Nil(t, err)

	err = T.m.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	data, err = ioutil.ReadFile("./tmp/data/test_data.txt")
	require.Equal(t, "this is changed data", string(data))

	err = T.m.DeleteFile("/test/data/test_data.txt")
	require.Nil(t, err)

	err = T.m.SyncRemoteToLocal("/test/data", "./tmp/data")
	require.Nil(t, err)

	_, err = os.Stat("./tmp/data/test_data.txt")
	require.NotNil(t, err)
}

func TestSeaweedFsManager_UpdateFile(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data", "/test/data")
	require.Nil(t, err)

	err = T.m.UpdateFile("/test/data/test_data.txt", []byte("this is changed data"))
	require.Nil(t, err)

	data, err := T.m.GetFile("/test/data/test_data.txt")
	require.Nil(t, err)
	require.Equal(t, "this is changed data", string(data))
}

func TestSeaweedFsManager_Exists(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.m.UploadDir("./data", "/test/data")
	require.Nil(t, err)

	ok, err := T.m.Exists("/test/data/test_data.txt")
	require.Nil(t, err)
	require.True(t, ok)

	ok, err = T.m.Exists("/test/data/test_data_404.txt")
	require.Nil(t, err)
	require.False(t, ok)
}
