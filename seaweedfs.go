package fs

import (
	"errors"
	"github.com/linxGnu/goseaweedfs"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

func NewSeaweedFSManager() (manager *SeaweedFSManager, err error) {
	filerUrl := viper.GetString("fs.seaweedfs.filer.url")
	if filerUrl == "" {
		filerUrl = "http://localhost:8888"
	}
	timeoutMinutes := viper.GetInt("fs.seaweed.filer.timeoutMinutes")
	if timeoutMinutes == 0 {
		timeoutMinutes = 5
	}
	filer, err := goseaweedfs.NewFiler(filerUrl, &http.Client{Timeout: time.Duration(timeoutMinutes) * time.Minute})
	if err != nil {
		return manager, err
	}
	manager = &SeaweedFSManager{
		f: filer,
	}
	if err := manager.Init(); err != nil {
		return manager, err
	}
	return
}

type SeaweedFSManager struct {
	f *goseaweedfs.Filer
}

func getCollectionAndTtlFromArgs(args ...interface{}) (collection, ttl string) {
	if len(args) > 0 {
		collection = args[0].(string)
	}
	if len(args) > 1 {
		ttl = args[1].(string)
	}
	return
}

func getUrlValuesFromArgs(args ...interface{}) (values url.Values) {
	if len(args) > 0 {
		values = args[0].(url.Values)
	}
	return values
}

func (s SeaweedFSManager) Init() (err error) {
	return nil
}

func (s SeaweedFSManager) Close() (err error) {
	if err := s.f.Close(); err != nil {
		return err
	}
	return nil
}

func (s SeaweedFSManager) ListDir(remotePath string, args ...interface{}) (files []goseaweedfs.FilerFileInfo, err error) {
	files, err = s.f.ListDirRecursive(remotePath)
	if err != nil {
		return files, err
	}
	return files, nil
}

func (s SeaweedFSManager) UploadFile(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	res, err := s.f.UploadFile(localPath, remotePath, collection, ttl)
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (s SeaweedFSManager) UploadDir(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	results, err := s.f.UploadDir(localPath, remotePath, collection, ttl)
	if err != nil {
		return err
	}
	for _, res := range results {
		if res.Error != "" {
			return errors.New(res.Error)
		}
	}
	return nil
}

func (s SeaweedFSManager) DownloadFile(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}
	urlValues := getUrlValuesFromArgs(args...)
	err = s.f.Download(remotePath, urlValues, func(reader io.Reader) error {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		dirPath := filepath.Dir(localPath)
		_, err = os.Stat(dirPath)
		if err != nil {
			// if not exists, create a new directory
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return err
			}
		}
		if err := ioutil.WriteFile(localPath, data, os.ModePerm); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s SeaweedFSManager) DownloadDir(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}
	files, err := s.ListDir(remotePath)
	for _, file := range files {
		if file.IsDir {
			if err := s.DownloadDir(file.FullPath, path.Join(localPath, file.Name), args...); err != nil {
				return err
			}
		} else {
			if err := s.DownloadFile(file.FullPath, path.Join(localPath, file.Name), args...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SeaweedFSManager) DeleteFile(remotePath string, args ...interface{}) (err error) {
	if err := s.f.DeleteFile(remotePath); err != nil {
		return err
	}
	return nil
}

func (s SeaweedFSManager) DeleteDir(remotePath string, args ...interface{}) (err error) {
	if err := s.f.DeleteDir(remotePath); err != nil {
		return err
	}
	return nil
}

func (s SeaweedFSManager) SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error) {
	// cache local files info
	localFilesMap := map[string]goseaweedfs.FileInfo{}
	localFiles, err := goseaweedfs.ListLocalFilesRecursive(localPath)
	if err != nil {
		return err
	}
	for _, file := range localFiles {
		localFilesMap[file.Path] = file
	}

	// cache remote files info
	remoteFilesMap := map[string]goseaweedfs.FilerFileInfo{}
	remoteFiles, err := s.f.ListDir(remotePath)
	if err != nil {
		return err
	}
	for _, file := range remoteFiles {
		remoteFilesMap[file.FullPath] = file
	}

	// TODO: compare both local and remote files

	return nil
}

func (s SeaweedFSManager) SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error) {
	panic("implement me")
}
