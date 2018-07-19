// +build dev

package assets

import (
	"go/build"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/union"
)

// Assets contains files that will be included in the binary
var Assets = union.New(map[string]http.FileSystem{
	"/migrations": http.Dir(importPathToDir("git.klink.asia/main/klinkregistry/assets/migrations")),
	"/static":     http.Dir(importPathToDir("git.klink.asia/main/klinkregistry/assets/static")),
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
