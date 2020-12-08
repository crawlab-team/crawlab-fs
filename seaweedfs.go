package fs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/linxGnu/goseaweedfs"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func getFilesAndFilesMaps(f *goseaweedfs.Filer, localPath, remotePath string) (localFiles []goseaweedfs.FileInfo, remoteFiles []goseaweedfs.FilerFileInfo, localFilesMap map[string]goseaweedfs.FileInfo, remoteFilesMap map[string]goseaweedfs.FilerFileInfo, err error) {
	// declare maps
	localFilesMap = map[string]goseaweedfs.FileInfo{}
	remoteFilesMap = map[string]goseaweedfs.FilerFileInfo{}

	// cache local files info
	localFiles, err = goseaweedfs.ListLocalFilesRecursive(localPath)
	if err != nil {
		return localFiles, remoteFiles, localFilesMap, remoteFilesMap, err
	}
	for _, file := range localFiles {
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(file.Path, localPath, "", -1))
		localFilesMap[fileRemotePath] = file
	}

	// cache remote files info
	remoteFiles, err = f.ListDir(remotePath)
	if err != nil {
		if err.Error() != FilerResponseNotFoundErrorMessage {
			return localFiles, remoteFiles, localFilesMap, remoteFilesMap, err
		}
		err = nil
	}
	for _, file := range remoteFiles {
		remoteFilesMap[file.FullPath] = file
	}

	return
}

func (s *SeaweedFSManager) Init() (err error) {
	return nil
}

func (s *SeaweedFSManager) Close() (err error) {
	if err := s.f.Close(); err != nil {
		return err
	}
	return nil
}

func (s *SeaweedFSManager) ListDir(remotePath string, args ...interface{}) (files []goseaweedfs.FilerFileInfo, err error) {
	files, err = s.f.ListDirRecursive(remotePath)
	if err != nil {
		return files, err
	}
	return files, nil
}

func (s *SeaweedFSManager) UploadFile(localPath, remotePath string, args ...interface{}) (err error) {
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

func (s *SeaweedFSManager) UploadDir(localPath, remotePath string, args ...interface{}) (err error) {
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

func (s *SeaweedFSManager) DownloadFile(remotePath, localPath string, args ...interface{}) (err error) {
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

func (s *SeaweedFSManager) DownloadDir(remotePath, localPath string, args ...interface{}) (err error) {
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

func (s *SeaweedFSManager) DeleteFile(remotePath string, args ...interface{}) (err error) {
	if err := s.f.DeleteFile(remotePath); err != nil {
		return err
	}
	return nil
}

func (s *SeaweedFSManager) DeleteDir(remotePath string, args ...interface{}) (err error) {
	if err := s.f.DeleteDir(remotePath); err != nil {
		return err
	}
	return nil
}

func (s *SeaweedFSManager) SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}

	// raise error if local path does not exist
	if _, err := os.Stat(localPath); err != nil {
		return err
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(s.f, localPath, remotePath)
	if err != nil {
		return err
	}

	// compare remote files with local files and delete files absent in local files
	for _, remoteFile := range remoteFiles {
		// skip directories
		if remoteFile.IsDir {
			continue
		}

		// attempt to get corresponding local file
		_, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, delete
			if err := s.DeleteFile(remoteFile.FullPath); err != nil {
				return err
			}
		}
	}

	// compare local files with remote files and upload files with difference
	for _, localFile := range localFiles {
		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(localFile.Path, localPath, "", -1))

		// attempt to get corresponding remote file
		remoteFile, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := s.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
				return err
			}
		} else {
			// file exists on remote, upload if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := s.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *SeaweedFSManager) SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}

	// create directory if local path does not exist
	if _, err := os.Stat(localPath); err != nil {
		if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
			return err
		}
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(s.f, localPath, remotePath)
	if err != nil {
		return err
	}

	// compare local files with remote files and delete files absent on remote
	for _, localFile := range localFiles {
		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(localFile.Path, localPath, "", -1))

		// attempt to get corresponding remote file
		_, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := os.Remove(localFile.Path); err != nil {
				return err
			}
		}
	}

	// compare remote files with local files and download if files with difference
	for _, remoteFile := range remoteFiles {
		// skip directories
		if remoteFile.IsDir {
			continue
		}

		// local file path
		localFileRelativePath := strings.Replace(remoteFile.FullPath, remotePath, "", -1)
		localFilePath := fmt.Sprintf("%s%s", localPath, localFileRelativePath)

		// attempt to get corresponding local file
		localFile, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, download
			if err := s.DownloadFile(remoteFile.FullPath, localFilePath); err != nil {
				return err
			}
		} else {
			// file exists on remote, download if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := s.DownloadFile(remoteFile.FullPath, localFilePath, args...); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *SeaweedFSManager) GetFile(remotePath string, args ...interface{}) (data []byte, err error) {
	urlValues := getUrlValuesFromArgs(args...)
	var buf bytes.Buffer
	err = s.f.Download(remotePath, urlValues, func(reader io.Reader) error {
		_, err := io.Copy(&buf, reader)
		if err != nil {
			return err
		}
		return nil
	})
	data = buf.Bytes()
	return
}

func (s *SeaweedFSManager) UpdateFile(remotePath string, data []byte, args ...interface{}) (err error) {
	tmpDirPath := "./tmp"
	if _, err := os.Stat(tmpDirPath); err != nil {
		if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
			return err
		}
	}
	tmpFilePath := path.Join(tmpDirPath, fmt.Sprintf(".%s", uuid.New().String()))
	if _, err := os.Stat(tmpFilePath); err == nil {
		if err := os.Remove(tmpFilePath); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(tmpFilePath, data, os.ModePerm); err != nil {
		return err
	}
	if err = s.UploadFile(tmpFilePath, remotePath, args...); err != nil {
		return err
	}
	if err := os.Remove(tmpFilePath); err != nil {
		return err
	}
	return
}

func (s *SeaweedFSManager) Exists(remotePath string, args ...interface{}) (ok bool, err error) {
	_, err = s.GetFile(remotePath, args...)
	if err == nil {
		// exists
		return true, nil
	}
	if strings.Contains(err.Error(), FilerStatusNotFoundErrorMessage) {
		// not exists
		return false, nil
	}
	return ok, err
}
