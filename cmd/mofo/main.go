package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/NicoNex/mofo"
)

// virtual file system serves a single file as if it were at the root "/".
type vFS string

func (vfs vFS) Open(name string) (http.File, error) {
	clean := path.Clean(name)

	// Accept "/" or the actual filename (e.g., "/document.md")
	if clean == "/" || clean == "." || clean == "/"+path.Base(string(vfs)) {
		return os.Open(string(vfs))
	}

	// Everything else → 404
	return nil, os.ErrNotExist
}

func retry(d time.Duration, fn func()) {
	for {
		fn()
		time.Sleep(d)
	}
}

func main() {
	var fs http.FileSystem

	port, path := parseFlags()
	if info, err := os.Stat(path); err != nil {
		log.Fatal(err)
	} else if info.IsDir() {
		fs = http.Dir(path)
	} else {
		fs = vFS(path)
	}

	http.Handle("/", mofo.FileServer(fs))
	retry(
		5*time.Second,
		func() { log.Println(http.ListenAndServe(port, nil)) },
	)
}

func parseFlags() (port, path string) {
	flag.StringVar(&port, "p", ":8080", "Port to serve on (shorthand)")
	flag.StringVar(&port, "port", ":8080", "Port to serve on")
	flag.Parse()

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	path = flag.Arg(0)
	if path == "" {
		flag.Usage()
		os.Exit(1)
	}
	return
}
