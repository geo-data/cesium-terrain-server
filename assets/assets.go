package assets

import (
	"fmt"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _data_gitignore = []byte(`# Ignore everything we don't need in this directory
smallterrain-blank.terrain
# Except this file
!.gitignore
`)

func data_gitignore_bytes() ([]byte, error) {
	return _data_gitignore, nil
}

func data_gitignore() (*asset, error) {
	bytes, err := data_gitignore_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "data/.gitignore", size: 110, mode: os.FileMode(420), modTime: time.Unix(1418210004, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _data_smallterrain_blank_terrain = []byte("\x1f\x8b\b\x00\x00\x00\x00\x00\x00\x03\xed\xce\xc1\r\x80 \x14\x04Qz\xb0h\x11\x9a\xf6@H0~\xaf\xbb\xdf0\x99\x02\xe6գ\n;\x83\x94\u007f\x04_\n\xbd\x80\x882t%\x10\xcc\xfc\x02\x97\x05\x81_\x10\xfdw\x13\xb4\x04\x82\x91_\xf0\x948\x05C\xe1\x164\x04\b6\x14\xf4\x04\x82Ƞ\x16\xf4\x97\xc3#X\x1dJ\x01\x11\x11\x11\x11\x11\x11\x11\xfd\xb3Rn\xd9\x11d\xac\x04!\x00\x00")

func data_smallterrain_blank_terrain_bytes() ([]byte, error) {
	return _data_smallterrain_blank_terrain, nil
}

func data_smallterrain_blank_terrain() (*asset, error) {
	bytes, err := data_smallterrain_blank_terrain_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "data/smallterrain-blank.terrain", size: 114, mode: os.FileMode(420), modTime: time.Unix(1418211210, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	"data/.gitignore": data_gitignore,
	"data/smallterrain-blank.terrain": data_smallterrain_blank_terrain,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"data": &_bintree_t{nil, map[string]*_bintree_t{
		".gitignore": &_bintree_t{data_gitignore, map[string]*_bintree_t{
		}},
		"smallterrain-blank.terrain": &_bintree_t{data_smallterrain_blank_terrain, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

