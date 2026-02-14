package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NicoNex/mofo"
)

func retry(d time.Duration, fn func()) {
	for {
		fn()
		time.Sleep(d)
	}
}

func main() {
	port, path := parseFlags()

	http.Handle("/", mofo.FileServer(http.Dir(path)))
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
