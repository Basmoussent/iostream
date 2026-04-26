package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web
var webFS embed.FS

// registerAssets serves the embedded SPA at /. Anything under /api stays
// reachable because http.ServeMux routes the more specific prefix first.
func registerAssets(mux *http.ServeMux) {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		// embed.FS guarantees this works at compile time; if it ever fails the
		// binary is corrupt and there is nothing useful to do at runtime.
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))
}
