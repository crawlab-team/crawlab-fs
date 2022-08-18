package fs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/crawlab-team/go-trace"
	"github.com/crawlab-team/goseaweedfs"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type SeaweedFsManager struct {
	// settings variables
	filerUrl string
	timeout  time.Duration
	authKey  string

	// internals
	f *goseaweedfs.Filer
}

func (m *SeaweedFsManager) Init() (err error) {
	var filerOpts []goseaweedfs.FilerOption
	if m.authKey != "" {
		filerOpts = append(filerOpts, goseaweedfs.WithFilerAuthKey(m.authKey))
	}
	m.f, err = goseaweedfs.NewFiler(m.filerUrl, &http.Client{Timeout: m.timeout}, filerOpts...)
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) Close() (err error) {
	if err := m.f.Close(); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) ListDir(remotePath string, isRecursive bool, args ...interface{}) (files []goseaweedfs.FilerFileInfo, err error) {
	if isRecursive {
		files, err = m.f.ListDirRecursive(remotePath)
	} else {
		files, err = m.f.ListDir(remotePath)
	}
	if err != nil {
		return files, trace.TraceError(err)
	}
	return files, nil
}

func (m *SeaweedFsManager) UploadFile(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	res, err := m.f.UploadFile(localPath, remotePath, collection, ttl)
	if err != nil {
		return trace.TraceError(err)
	}
	if res.Error != "" {
		err = errors.New(res.Error)
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) UploadDir(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	results, err := m.f.UploadDir(localPath, remotePath, collection, ttl)
	if err != nil {
		return trace.TraceError(err)
	}
	for _, res := range results {
		if res.Error != "" {
			err = errors.New(res.Error)
			return trace.TraceError(err)
		}
	}
	return nil
}

func (m *SeaweedFsManager) DownloadFile(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	urlValues := getUrlValuesFromArgs(args...)
	err = m.f.Download(remotePath, urlValues, func(reader io.Reader) error {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return trace.TraceError(err)
		}
		dirPath := filepath.Dir(localPath)
		_, err = os.Stat(dirPath)
		if err != nil {
			// if not exists, create a new directory
			if err := os.MkdirAll(dirPath, DefaultDirMode); err != nil {
				return trace.TraceError(err)
			}
		}
		fileMode := DefaultFileMode
		fileInfo, err := os.Stat(localPath)
		if err == nil {
			// if file already exists, save file mode and remove it
			fileMode = fileInfo.Mode()
			if err := os.Remove(localPath); err != nil {
				return trace.TraceError(err)
			}
		}
		if err := ioutil.WriteFile(localPath, data, fileMode); err != nil {
			return trace.TraceError(err)
		}
		return nil
	})
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) DownloadDir(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	files, err := m.ListDir(remotePath, true)
	for _, file := range files {
		if file.IsDir {
			if err := m.DownloadDir(file.FullPath, path.Join(localPath, file.Name), args...); err != nil {
				return trace.TraceError(err)
			}
		} else {
			if err := m.DownloadFile(file.FullPath, path.Join(localPath, file.Name), args...); err != nil {
				return trace.TraceError(err)
			}
		}
	}
	return nil
}

func (m *SeaweedFsManager) DeleteFile(remotePath string, args ...interface{}) (err error) {
	if err := m.f.DeleteFile(remotePath); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) DeleteDir(remotePath string, args ...interface{}) (err error) {
	if err := m.f.DeleteDir(remotePath); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}

	// raise error if local path does not exist
	if _, err := os.Stat(localPath); err != nil {
		return trace.TraceError(err)
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(m.f, localPath, remotePath)
	if err != nil {
		return trace.TraceError(err)
	}

	// compare remote files with local files and delete files absent in local files
	for _, remoteFile := range remoteFiles {
		// attempt to get corresponding local file
		_, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, delete
			if remoteFile.IsDir {
				if err := m.DeleteDir(remoteFile.FullPath); err != nil {
					return trace.TraceError(err)
				}
			} else {
				if err := m.DeleteFile(remoteFile.FullPath); err != nil {
					return trace.TraceError(err)
				}
			}
		}
	}

	// compare local files with remote files and upload files with difference
	for _, localFile := range localFiles {
		// skip .git
		if IsGitFile(localFile) {
			continue
		}

		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(localFile.Path, localPath, "", -1))

		// attempt to get corresponding remote file
		remoteFile, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := m.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
				return trace.TraceError(err)
			}
		} else {
			// file exists on remote, upload if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := m.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
					return trace.TraceError(err)
				}
			}
		}
	}

	return nil
}

func (m *SeaweedFsManager) SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error) {
	// create directory if local path does not exist
	if _, err := os.Stat(localPath); err != nil {
		if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
			return trace.TraceError(err)
		}
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(m.f, localPath, remotePath)
	if err != nil {
		return trace.TraceError(err)
	}

	// compare local files with remote files and delete files absent on remote
	for _, localFile := range localFiles {
		// skip .git
		if IsGitFile(localFile) {
			continue
		}

		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(localFile.Path, localPath, "", -1))

		// attempt to get corresponding remote file
		_, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := os.Remove(localFile.Path); err != nil {
				return trace.TraceError(err)
			}
		}
	}

	// compare remote files with local files and download if files with difference
	for _, remoteFile := range remoteFiles {
		// directory
		if remoteFile.IsDir {
			localDirRelativePath := strings.Replace(remoteFile.FullPath, remotePath, "", 1)
			localDirPath := fmt.Sprintf("%s%s", localPath, localDirRelativePath)
			if err := m.SyncRemoteToLocal(remoteFile.FullPath, localDirPath); err != nil {
				return err
			}
			continue
		}

		// local file path
		localFileRelativePath := strings.Replace(remoteFile.FullPath, remotePath, "", 1)
		localFilePath := fmt.Sprintf("%s%s", localPath, localFileRelativePath)

		// attempt to get corresponding local file
		localFile, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, download
			if err := m.DownloadFile(remoteFile.FullPath, localFilePath); err != nil {
				return trace.TraceError(err)
			}
		} else {
			// file exists on remote, download if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := m.DownloadFile(remoteFile.FullPath, localFilePath, args...); err != nil {
					return trace.TraceError(err)
				}
			}
		}
	}

	return nil
}

func (m *SeaweedFsManager) GetFile(remotePath string, args ...interface{}) (data []byte, err error) {
	urlValues := getUrlValuesFromArgs(args...)
	var buf bytes.Buffer
	err = m.f.Download(remotePath, urlValues, func(reader io.Reader) error {
		_, err := io.Copy(&buf, reader)
		if err != nil {
			return trace.TraceError(err)
		}
		return nil
	})
	data = buf.Bytes()
	return
}

func (m *SeaweedFsManager) GetFileInfo(remotePath string, args ...interface{}) (file *goseaweedfs.FilerFileInfo, err error) {
	arr := strings.Split(remotePath, "/")
	dirName := strings.Join(arr[:(len(arr)-1)], "/")
	files, err := m.f.ListDir(dirName)
	if err != nil {
		return file, trace.TraceError(err)
	}
	for _, f := range files {
		if f.FullPath == remotePath {
			return &f, nil
		}
	}
	return nil, trace.TraceError(ErrorFsNotExists)
}

func (m *SeaweedFsManager) UpdateFile(remotePath string, data []byte, args ...interface{}) (err error) {
	tmpRootDir := os.TempDir()
	tmpDirPath := path.Join(tmpRootDir, ".seaweedfs")
	if _, err := os.Stat(tmpDirPath); err != nil {
		if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
			return trace.TraceError(err)
		}
	}
	tmpFilePath := path.Join(tmpDirPath, fmt.Sprintf(".%s", uuid.New().String()))
	if _, err := os.Stat(tmpFilePath); err == nil {
		if err := os.Remove(tmpFilePath); err != nil {
			return trace.TraceError(err)
		}
	}
	if err := ioutil.WriteFile(tmpFilePath, data, os.ModePerm); err != nil {
		return trace.TraceError(err)
	}
	if err = m.UploadFile(tmpFilePath, remotePath, args...); err != nil {
		return trace.TraceError(err)
	}
	if err := os.Remove(tmpFilePath); err != nil {
		return trace.TraceError(err)
	}
	return
}

func (m *SeaweedFsManager) Exists(remotePath string, args ...interface{}) (ok bool, err error) {
	_, err = m.GetFile(remotePath, args...)
	if err == nil {
		// exists
		return true, nil
	}
	if strings.Contains(err.Error(), FilerStatusNotFoundErrorMessage) {
		// not exists
		return false, nil
	}
	return ok, trace.TraceError(err)
}

func (m *SeaweedFsManager) SetFilerUrl(url string) {
	m.filerUrl = url
}

func (m *SeaweedFsManager) SetFilerAuthKey(authKey string) {
	m.authKey = authKey
}

func (m *SeaweedFsManager) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func NewSeaweedFsManager(opts ...Option) (m2 Manager, err error) {
	// manager
	m := &SeaweedFsManager{
		filerUrl: "http://localhost:8888",
		timeout:  5 * time.Minute,
	}

	// apply options
	for _, opt := range opts {
		opt(m)
	}

	// initialize
	if err := m.Init(); err != nil {
		return nil, err
	}

	return m, nil
}
