package fs

import "github.com/linxGnu/goseaweedfs"

type Manager interface {
	Init() (err error)
	Close() (err error)
	ListDir(remotePath string, isRecursive bool, args ...interface{}) (files []goseaweedfs.FilerFileInfo, err error)
	UploadFile(localPath, remotePath string, args ...interface{}) (err error)
	UploadDir(localPath, remotePath string, args ...interface{}) (err error)
	DownloadFile(remotePath, localPath string, args ...interface{}) (err error)
	DownloadDir(remotePath, localPath string, args ...interface{}) (err error)
	DeleteFile(remotePath string, args ...interface{}) (err error)
	DeleteDir(remotePath string, args ...interface{}) (err error)
	SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error)
	SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error)
	GetFile(remotePath string, args ...interface{}) (data []byte, err error)
	UpdateFile(remotePath string, data []byte, args ...interface{}) (err error)
	Exists(remotePath string, args ...interface{}) (ok bool, err error)
}
