// Command sbi serves the SBI feedback app.
//
// The app itself is static: everything under static/ is embedded into the
// binary and no request data is read or stored. Asset paths in the HTML are
// relative so the app works at / locally and under /sbi/ behind a proxy.
package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFS embed.FS

func handler() http.Handler {
	app, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(app))
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	log.Printf("sbi listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, handler()))
}
