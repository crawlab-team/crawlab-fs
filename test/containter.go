// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// docker-compose.yml
package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dockerComposeYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\xc1\x6e\x83\x30\x0c\x86\xef\x3c\x85\x45\x0f\x9c\x50\xda\xf5\xc2\x22\xf5\x31\x76\xae\x3c\x30\x6d\x26\x92\x20\x3b\x50\x75\x4f\x3f\x25\x41\xd0\xed\x30\xf5\xd2\x8b\x95\xfc\xb2\x3f\xc7\xbf\x33\x13\x8b\xf1\x4e\x43\xf5\x56\x15\x85\x10\xcf\xa6\x25\xd1\x05\x80\x45\x09\xc4\xf1\x04\x60\x2c\x5e\x48\x43\x7b\x65\x23\xc3\x24\xbd\x12\xc2\x1b\x51\xd7\x0b\xec\x60\x12\x02\x04\x26\xeb\x03\xe5\xcc\x54\xd3\x7a\x17\xd0\x38\xe2\xb3\x43\x4b\x1a\xd6\x92\x3a\x93\x53\x12\x93\x04\xe4\xa0\x01\x87\x1b\xde\x25\x69\xa3\xe7\x20\xb9\x2f\x40\x0d\xef\xc7\xe3\x51\xc7\xb0\x2a\x87\x24\x1d\x56\xad\xf5\xd6\xa2\xeb\x34\x94\x99\x0c\xb5\x19\x4f\xf9\x58\xa6\x84\xd9\x0f\x93\xa5\x07\x66\xa9\x3e\x84\x58\x94\x45\x9e\xbf\xaf\xe8\x2e\x6a\x64\xff\x45\x6d\x10\xd5\x32\xde\x06\xfc\xac\x03\xa1\x5d\x2f\xbd\xa8\x60\x47\xb5\x58\xa2\xd8\xfb\xa0\xb6\x06\x19\xff\x0a\xa7\x32\xf9\x59\xa7\x9a\x7d\xb3\xd7\x31\x6c\x4e\x25\xe9\xb0\x6a\xab\x53\x55\x26\x43\x6d\xe3\xca\x89\x4f\x8b\x75\xc9\xe8\x12\xea\x48\x3e\xc5\xb2\x2a\xd5\x75\x34\x92\xeb\xe4\xec\xdd\xd6\x6d\xdd\x62\x6f\x86\xd7\xfc\x93\x04\x7e\x7a\xf8\xa6\x69\x74\x0c\xc5\x6e\x9b\x3e\x6a\x29\xfe\x99\x3e\x91\x61\xf9\x88\xbf\x67\xcf\x03\x87\x70\xd7\x10\x78\xca\x4f\x94\xd0\x19\x77\xf6\x23\xb9\x07\xf1\x5f\x4f\xf2\x75\x59\xdf\x4f\x00\x00\x00\xff\xff\xff\xb5\xe4\x90\x64\x03\x00\x00")

func dockerComposeYmlBytes() ([]byte, error) {
	return bindataRead(
		_dockerComposeYml,
		"docker-compose.yml",
	)
}

func dockerComposeYml() (*asset, error) {
	bytes, err := dockerComposeYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "docker-compose.yml", size: 868, mode: os.FileMode(420), modTime: time.Unix(1608801983, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"docker-compose.yml": dockerComposeYml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"docker-compose.yml": &bintree{dockerComposeYml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
