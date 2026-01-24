package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

type Server struct {
	Root       string
	MD         goldmark.Markdown
	HTTPServer *http.Server
}

type PageData struct {
	Title       string
	Description string
	Year        int
	Content     template.HTML
}

type Frontmatter struct {
	Title       string `yaml:"title" toml:"title"`
	Description string `yaml:"description" toml:"description"`
}

var (
	//go:embed page.template.html
	tmplContent string
	//go:embed default.style.css
	defaultCSS string
)

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

func newServer(port string) Server {
	return Server{
		// Port: port,
		MD: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Typographer,
				&frontmatter.Extender{
					Formats: frontmatter.DefaultFormats,
				},
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
		HTTPServer: &http.Server{
			Addr: port,
		},
	}
}

func (s *Server) Serve(path string) error {
	clean := filepath.Clean(path)

	info, err := os.Stat(clean)
	if err != nil {
		return err
	}

	var file string
	if info.IsDir() {
		idx := filepath.Join(clean, "index.md")
		if _, err := os.Stat(idx); errors.Is(err, fs.ErrNotExist) {
			log.Fatal("Serve", "os.Stat", err)
		}
		file = idx
		s.Root = clean
	} else {
		file = clean
		s.Root = filepath.Dir(clean)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.serveFile(w, r, file)
	})

	s.HTTPServer.Handler = mux
	return s.HTTPServer.ListenAndServe()
}

func (s Server) serveMarkdown(w http.ResponseWriter, _ *http.Request, mdPath string) {
	content, err := os.ReadFile(mdPath)
	if err != nil {
		log.Println("serveMarkdown", "os.ReadFile", err)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	ctx := parser.NewContext()

	var buf strings.Builder
	if err := s.MD.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		log.Println("serveMarkdown", "MD.Convert", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	fm := Frontmatter{Title: "Untitled"}
	if data := frontmatter.Get(ctx); data != nil {
		if err := data.Decode(&fm); err != nil {
			log.Println("serveMarkdown", "frontmatter.Decode", err)
		}
	}

	if fm.Title == "" {
		fm.Title = "Untitled"
	}

	tmpl, err := template.New("page").Parse(tmplContent)
	if err != nil {
		log.Println("serveMarkdown", "template.Parse", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	vm := PageData{
		Title:       fm.Title,
		Description: fm.Description,
		Year:        time.Now().Year(),
		Content:     template.HTML(buf.String()),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, vm); err != nil {
		log.Println("serveMarkdown", "template.Execute", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func (s Server) serveCSS(w http.ResponseWriter, r *http.Request) {
	css := filepath.Join(s.Root, "assets", "style.css")
	if _, err := os.Stat(css); err == nil {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, css)
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write([]byte(defaultCSS))
}

func (s Server) serveFile(w http.ResponseWriter, r *http.Request, file string) {
	if r.URL.Path == "/" {
		s.serveMarkdown(w, r, file)
		return
	}

	if r.URL.Path == "/assets/style.css" {
		s.serveCSS(w, r)
		return
	}

	trimmed := strings.TrimPrefix(r.URL.Path, "/")
	safe, err := validatePath(s.Root, trimmed)
	if err != nil {
		log.Println("serveFile", "validatePath", err)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if strings.HasSuffix(strings.ToLower(safe), ".md") {
		s.serveMarkdown(w, r, safe)
		return
	}

	http.ServeFile(w, r, safe)
}

func validatePath(root, p string) (string, error) {
	if !fs.ValidPath(p) {
		return "", fmt.Errorf("invalid path: %s", p)
	}

	full := filepath.Join(root, p)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		return "", fmt.Errorf("not found: %s", p)
	}

	return full, nil
}

func retry(d time.Duration, fn func()) {
	for {
		fn()
		time.Sleep(d)
	}
}

func main() {
	port, path := parseFlags()
	server := newServer(port)

	retry(
		5*time.Second,
		func() { log.Println(server.Serve(path)) },
	)
}
