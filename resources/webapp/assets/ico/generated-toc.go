// AUTO-GENERATED TOC
// DO NOT EDIT!!!
package webapp_assets_ico

import (
	"net/http"
	"os"
	embedfs "../resources"
)

func init() {

	DIR.AddDir(DIR)

}

var DIR = embedfs.DirAlloc("ico")

func Dir(path string) http.FileSystem {
	if handle, err := DIR.Open(); err == nil {
		return handle
	}
	return nil
}

func Mount() http.FileSystem {
	return Dir(".")
}

func FileInfo() os.FileInfo {
	return DIR
}
