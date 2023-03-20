package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

var dist embed.FS

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

// Handler delivers the application frontend
func Handler(prefix, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		// If we can't find the asset, return the default index.html content
		f, err := dist.Open(assetPath)
		if os.IsNotExist(err) {
			return dist.Open("dist/index.html")
		}
		// Otherwise assume this is a legitimate request routed correctly
		return f, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}
