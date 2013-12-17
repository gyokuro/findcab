// AUTO-GENERATED TOC
// DO NOT EDIT!!!
package webapp

import (
	"net/http"
	"os"
	embedfs "../resources"
)

import (
	assets "../resources/webapp/assets"

	dist "../resources/webapp/dist"

	fonts "../resources/webapp/fonts"
)

func init() {

	DIR.AddDir(DIR)

	DIR.AddDir(assets.DIR)

	DIR.AddDir(dist.DIR)

	DIR.AddDir(fonts.DIR)

}

var DIR = embedfs.DirAlloc("webapp")

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
