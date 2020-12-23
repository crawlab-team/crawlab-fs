package fs

import (
	"github.com/linxGnu/goseaweedfs"
	"regexp"
)

func IsGitFile(file goseaweedfs.FileInfo) (res bool) {
	// skip .git
	res, err := regexp.MatchString("/?\\.git/", file.Path)
	if err != nil {
		return false
	}
	return res
}
