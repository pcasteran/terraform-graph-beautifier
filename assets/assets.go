// +build dev

package assets

//go:generate go run -tags=dev assets_generate.go

import (
	"net/http"
	"path"
	"runtime"
)

// resolvePath resolves the specified path according to this module directory instead of the current directory as returned by `os.Getwd()`.
func resolvePath(p string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	moduleDir := path.Dir(filename)
	return path.Join(moduleDir, p)
}

// Templates contains the HTML templates.
var Templates http.FileSystem = http.Dir(resolvePath("../templates"))
