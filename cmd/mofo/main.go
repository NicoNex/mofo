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

// virtual file system serves a single file as if it were in a root directory "/".
type vFS string

func (vfs vFS) Open(name string) (http.File, error) {
	clean := path.Clean(name)

	// Serve the file as /index.md so mofo.FileServer serves it at /
	// Also serve at the original basename for direct access
	if clean == "/index.html" || clean == "/"+path.Base(string(vfs)) {
		return os.Open(string(vfs))
	}

	// Root directory request - return virtual directory
	if clean == "/" || clean == "." {
		realFile, err := os.Open(string(vfs))
		if err != nil {
			return nil, err
		}
		stat, err := realFile.Stat()
		realFile.Close()
		if err != nil {
			return nil, err
		}
		return &vDir{name: path.Base(string(vfs)), info: stat}, nil
	}

	// Everything else → 404
	return nil, os.ErrNotExist
}

// vDir implements http.File for a virtual directory containing one file.
type vDir struct {
	name   string
	info   os.FileInfo
	offset int
}

func (d *vDir) Close() error                   { return nil }
func (d *vDir) Read([]byte) (int, error)       { return 0, os.ErrInvalid }
func (d *vDir) Seek(int64, int) (int64, error) { return 0, os.ErrInvalid }

func (d *vDir) Stat() (os.FileInfo, error) {
	return dirInfo("/"), nil
}

func (d *vDir) Readdir(count int) ([]os.FileInfo, error) {
	if d.offset > 0 {
		return nil, nil
	}
	d.offset++
	return []os.FileInfo{d.info}, nil
}

// dirInfo implements os.FileInfo for the virtual root directory.
type dirInfo string

func (di dirInfo) Name() string       { return string(di) }
func (di dirInfo) Size() int64        { return 0 }
func (di dirInfo) Mode() os.FileMode  { return os.ModeDir | 0755 }
func (di dirInfo) ModTime() time.Time { return time.Now() }
func (di dirInfo) IsDir() bool        { return true }
func (di dirInfo) Sys() any           { return nil }

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
