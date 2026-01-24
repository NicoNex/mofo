# mofo

A dead-simple markdown web server built with the [motherfucking website](https://motherfuckingwebsite.com/) philosophy in mind: fast, minimal, and brutally functional.

## Philosophy

This tool embraces simplicity:
- **No JavaScript** - Just pure HTML and CSS
- **Minimal CSS** - Clean, readable typography with no bloat
- **Fast** - Serves markdown files with zero build step
- **Portable** - Single binary, no dependencies

Perfect for personal sites, documentation, blogs, or any content that values substance over style.

## Features

- 🔥 Serves markdown files as beautifully rendered HTML
- 📝 Supports frontmatter (YAML/TOML) for metadata
- 🎨 Uses a minimal default stylesheet (customizable)
- 🚀 Hot-reload friendly (auto-retries on crash)
- 📁 Serves static assets from the same directory
- 🔗 GitHub Flavored Markdown support
- 📖 Auto-generates heading IDs for anchor links

## Installation

```bash
go install github.com/yourusername/mofo@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/mofo.git
cd mofo
go build
```

## Usage

### Serve a single markdown file

```bash
mofo /path/to/file.md
```

### Serve a directory (requires index.md)

```bash
mofo /path/to/directory
```

The directory must contain an `index.md` file as the entry point.

### Specify a custom port

```bash
mofo -p 3000 /path/to/directory
mofo -port 3000 /path/to/directory
```

Default port is `:8080`.

## Directory Structure

```
your-site/
├── index.md              # Entry point (required)
├── about.md              # Other markdown pages
├── posts/
│   └── hello-world.md
└── assets/               # Static assets folder
    ├── style.css         # Custom stylesheet (overrides default)
    ├── image.png         # Images
    ├── script.js         # Scripts (if you must...)
    └── logo.svg          # Any other static files
```

## Assets Folder

The `assets/` folder works exactly like a standard website's static files directory. Any file placed in `assets/` can be referenced from your markdown or HTML:

```markdown
![My Image](/assets/image.png)
<img src="/assets/logo.svg" alt="Logo">
<script src="/assets/script.js"></script>
```

All files in the `assets/` directory are served directly with no processing.

## Custom Styling

**To override the default CSS:** Simply create a file at `assets/style.css` in your site's root directory.

```
your-site/
└── assets/
    └── style.css  # Overrides the built-in minimal CSS
```

- If `assets/style.css` **exists** → your custom CSS is used
- If `assets/style.css` **does not exist** → the built-in minimal CSS is used

The default CSS is embedded in the binary, so you always have a working baseline. Create your own `style.css` only when you want to customize the look.

## Frontmatter

Add metadata to your markdown files using frontmatter:

```markdown
---
title: My Awesome Page
meta: A short description for SEO
---

# Content starts here

Your markdown content...
```

**Supported formats:** YAML and TOML

**Supported fields:**
- `title` - Page title (shown in browser tab and `<title>` tag)
- `meta` - Meta description (currently parsed but not yet rendered)

If no title is provided, defaults to "Untitled".

## URL Structure

- `/` - Serves the index.md file
- `/about.md` - Serves about.md from the root
- `/posts/hello-world.md` - Serves posts/hello-world.md
- `/assets/style.css` - Serves the stylesheet (custom or default)
- `/assets/image.png` - Serves static assets directly
- `/assets/anything.xyz` - Serves any file from assets/

All paths are relative to the root directory you specified when starting the server.

## Command Line Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `-port` | `-p` | `:8080` | Port to serve on |

**Important:** Flags must come **before** the path argument.

```bash
# ✅ Correct
mofo -p 3000 /path/to/site

# ❌ Wrong
mofo /path/to/site -p 3000
```

## Examples

### Minimal personal site

```bash
# Create a simple site
mkdir my-site
cd my-site
echo "---\ntitle: Home\n---\n\n# Welcome\n\nThis is my site." > index.md

# Serve it
mofo .
```

Visit `http://localhost:8080` in your browser.

### Blog with custom styling

```bash
my-blog/
├── index.md
├── posts/
│   ├── 2024-01-15-first-post.md
│   └── 2024-01-20-second-post.md
└── assets/
    ├── style.css         # Your custom CSS
    └── header-image.jpg  # Referenced in markdown
```

Serve with: `mofo my-blog`

Reference the image in your markdown:
```markdown
![Header](/assets/header-image.jpg)
```

## How It Works

1. **Starts** an HTTP server on the specified port
2. **Reads** markdown files from the specified directory
3. **Parses** frontmatter and converts markdown to HTML using [goldmark](https://github.com/yuin/goldmark)
4. **Renders** the HTML using a minimal embedded template
5. **Serves** static assets (images, CSS, etc.) directly from the `assets/` folder
6. **Retries** automatically if the server crashes (every 5 seconds)

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [goldmark/frontmatter](https://go.abhg.dev/goldmark/frontmatter) - Frontmatter support

## License

MIT

## Contributing

Contributions welcome! Keep it simple, keep it fast, keep it minimal.

---

*Built with the motherfucking website philosophy: content first, bullshit never.*
