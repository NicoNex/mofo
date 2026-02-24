package mofo

import (
	_ "embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

// Frontmatter data.
type pageMeta struct {
	Title       string `yaml:"title" toml:"title"`
	Description string `yaml:"description" toml:"description"`
}

type pageData struct {
	Title       string
	Description string
	Year        int
	Content     template.HTML
	Meta        pageMeta
}

type fileServer struct {
	handler  http.Handler
	md       goldmark.Markdown
	root     http.FileSystem
	template *template.Template
}

var (
	//go:embed page.template.html
	HTMLtemplate string
	//go:embed default.style.css
	DefaultCSS []byte
)

func FileServer(root http.FileSystem) http.Handler {
	return http.HandlerFunc(
		fileServer{
			// Set the root directory
			root: root,
			// Standard file server for regular files
			handler: http.FileServer(root),
			// Initialize HTML template for markdown pages
			template: template.Must(template.New("mdpage").Parse(HTMLtemplate)),
			// Initialize markdown processor
			md: goldmark.New(
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
		}.handleRequest,
	)
}

func (fs fileServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Reject path traversal attempts with backslash (Windows-style)
	if strings.Contains(r.URL.Path, "\\") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch path := r.URL.Path; {
	case path == "/":
		if err := fs.serveMarkdown(w, "/index.md"); err == nil {
			return
		}

	case path == "/assets/style.css":
		fs.serveCSS(w, r)
		return

	case strings.HasSuffix(strings.ToLower(path), ".md"):
		if fs.serveMarkdown(w, path) == nil {
			return
		}
	}

	// Delegate to standard file server for everything else
	fs.handler.ServeHTTP(w, r)
}

func (fs fileServer) serveMarkdown(w http.ResponseWriter, mdPath string) error {
	html, meta, err := fs.convertMarkdown(mdPath)
	if err != nil {
		return err
	}

	fs.renderPage(w, html, meta)
	return nil
}

func (fs fileServer) convertMarkdown(mdPath string) (string, pageMeta, error) {
	f, err := fs.root.Open(mdPath)
	if err != nil {
		return "", pageMeta{}, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", pageMeta{}, err
	}

	ctx := parser.NewContext()

	var buf strings.Builder
	if err := fs.md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return "", pageMeta{}, err
	}

	meta := pageMeta{Title: "Untitled"}
	if data := frontmatter.Get(ctx); data != nil {
		if err := data.Decode(&meta); err != nil {
			log.Println("convertMarkdown", "frontmatter.Decode", err)
		}
	}

	return buf.String(), meta, nil
}

func (fs fileServer) renderPage(w http.ResponseWriter, html string, meta pageMeta) {
	vm := pageData{
		Title:       meta.Title,
		Description: meta.Description,
		Year:        time.Now().Year(),
		Content:     template.HTML(html),
		Meta:        meta,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := fs.template.Execute(w, vm); err != nil {
		log.Println("renderPage", "template.Execute", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func (fs fileServer) serveCSS(w http.ResponseWriter, r *http.Request) {
	// Try to serve custom CSS
	if f, err := fs.root.Open("/assets/style.css"); err == nil {
		defer f.Close()
		if stat, err := f.Stat(); err == nil {
			http.ServeContent(w, r, "style.css", stat.ModTime(), f)
			return
		}
	}

	// Serve default CSS
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(DefaultCSS)
}
