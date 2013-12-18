// AUTO-GENERATED TOC
// DO NOT EDIT!!!
package webapp_dist

import (
	"net/http"
	"os"
	embedfs "github.com/gyokuro/findcab/resources"
)

import (
	css "github.com/gyokuro/findcab/resources/webapp/dist/css"

	js "github.com/gyokuro/findcab/resources/webapp/dist/js"
)

func init() {

	DIR.AddDir(DIR)

	DIR.AddDir(css.DIR)

	DIR.AddDir(js.DIR)

}

var DIR = embedfs.DirAlloc("dist")

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
