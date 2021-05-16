// Code generated for package test by go-bindata DO NOT EDIT. (@generated)
// sources:
// bin/start.sh
// bin/stop.sh
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

var _binStartSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\xc9\x41\xaa\x02\x31\x10\x45\xd1\x79\xad\xe2\x7d\x3e\x38\xeb\x64\xae\x4b\x51\x07\xad\x79\x21\x85\x49\xa7\x49\xc5\xd6\xe5\x8b\x04\x44\x1c\xde\x7b\xfe\xff\xfc\x45\x17\x6f\x49\x34\xe2\x88\x89\x70\xbe\x97\x15\xe7\x03\x7a\xe2\x22\xc0\x5e\x98\x8d\x02\x94\x5b\xd0\x36\x58\xa2\xca\x83\x0c\x30\xb6\x8d\x0d\x27\x01\xa6\x8f\x8e\x2c\xb3\x75\x36\xf7\x73\xb7\x9a\xef\x85\xef\xeb\x34\x3c\xbf\x45\x57\xe4\x7a\x9d\x73\xaa\xd6\xc7\x89\x9a\xd9\xb0\x93\x57\x00\x00\x00\xff\xff\x31\x83\x37\x15\xa4\x00\x00\x00")

func binStartShBytes() ([]byte, error) {
	return bindataRead(
		_binStartSh,
		"bin/start.sh",
	)
}

func binStartSh() (*asset, error) {
	bytes, err := binStartShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bin/start.sh", size: 164, mode: os.FileMode(493), modTime: time.Unix(1621151636, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _binStopSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\xd4\x4f\xca\xcc\xd3\x2f\xce\xe0\xe2\xca\xce\xcc\xc9\x51\xd0\xb5\x54\x48\x28\x28\x56\x48\xac\x28\xad\x49\x2f\x4a\x2d\x50\x28\x4f\x4d\x4d\x81\xb0\x74\xcb\x14\x40\x74\x4d\x62\x79\xb6\x82\x7a\x75\x41\x51\x66\x5e\x89\x82\x8a\x51\xad\x7a\x4d\x45\x62\x51\x7a\x71\x02\x20\x00\x00\xff\xff\xe9\x34\x7d\x31\x49\x00\x00\x00")

func binStopShBytes() ([]byte, error) {
	return bindataRead(
		_binStopSh,
		"bin/stop.sh",
	)
}

func binStopSh() (*asset, error) {
	bytes, err := binStopShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bin/stop.sh", size: 73, mode: os.FileMode(493), modTime: time.Unix(1620995602, 0)}
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
	"bin/start.sh": binStartSh,
	"bin/stop.sh":  binStopSh,
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
	"bin": &bintree{nil, map[string]*bintree{
		"start.sh": &bintree{binStartSh, map[string]*bintree{}},
		"stop.sh":  &bintree{binStopSh, map[string]*bintree{}},
	}},
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
