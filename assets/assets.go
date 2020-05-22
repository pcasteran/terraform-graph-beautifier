// +build dev

package assets

//go:generate go run -tags=dev assets_generate.go

import "net/http"

// Templates contains the HTML templates.
var Templates http.FileSystem = http.Dir("../templates")
