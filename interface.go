package fs

import "github.com/linxGnu/goseaweedfs"

type Manager interface {
	Init() (err error)
	Close() (err error)
	ListDir(remotePath string, args ...interface{}) (files []goseaweedfs.FilerFileInfo, err error)
	UploadFile(localPath, remotePath string, args ...interface{}) (err error)
	UploadDir(localPath, remotePath string, args ...interface{}) (err error)
	DownloadFile(remotePath, localPath string, args ...interface{}) (err error)
	DownloadDir(remotePath, localPath string, args ...interface{}) (err error)
	DeleteFile(remotePath string, args ...interface{}) (err error)
	DeleteDir(remotePath string, args ...interface{}) (err error)
	SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error)
	SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error)
}
