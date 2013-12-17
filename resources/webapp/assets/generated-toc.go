// AUTO-GENERATED TOC
// DO NOT EDIT!!!
package webapp_assets

import (
	"net/http"
	"os"
	embedfs "../resources"
)

import (
	css "../resources/webapp/assets/css"

	ico "../resources/webapp/assets/ico"

	js "../resources/webapp/assets/js"
)

func init() {

	DIR.AddDir(DIR)

	DIR.AddDir(css.DIR)

	DIR.AddDir(ico.DIR)

	DIR.AddDir(js.DIR)

}

var DIR = embedfs.DirAlloc("assets")

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
