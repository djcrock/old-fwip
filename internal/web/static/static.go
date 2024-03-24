package static

import (
	"embed"
	"net/http"
)

//go:embed *.css *.html *.js
var files embed.FS

var FileServer = http.FileServer(http.FS(files))

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, files, "index.html")
}
