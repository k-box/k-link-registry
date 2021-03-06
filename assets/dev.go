// +build dev

package assets

import (
	"go/build"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/union"
)

// Assets contains files that will be included in the binary
// this is a union file system, so to reach the expected file
// the root folder defined in the map should be prepended
// to the file path
var Assets = union.New(map[string]http.FileSystem{
	"/migrations": http.Dir(importPathToDir("github.com/k-box/k-link-registry/assets/migrations")),
	"/static":     http.Dir(importPathToDir("github.com/k-box/k-link-registry/ui/dist")),
})

// importPathToDir is a helper function that resolves the absolute path of
// modules, so they can be used both in dev mode (`-tags="dev"`) or with a
// generated static asset file (`go generate`).
func importPathToDir(importPath string) string {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		log.Fatalln(err)
	}
	return p.Dir
}
