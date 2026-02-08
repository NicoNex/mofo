# mofo

## This is a motherfucking markdown web server.

And it's probably the only one you'll ever fucking need.

### It's *actually* simple

You know what this does? It serves markdown files as HTML and everything else as static files. **That's it.** No build step. No npm install bullshit that downloads half the internet. No framework that's "totally different this time, bro." Just a single goddamn binary that does exactly what it says on the tin.

Files ending in `.md`? Rendered as HTML. Everything else? Served directly. Images, PDFs, CSS, JavaScript (if you hate yourself), fonts, videos, carrier pigeon instructions - doesn't matter. If it's not markdown, it gets served as-is.

### Seriously, it's literally one command

```bash
mofo /path/to/your/site
```

**BOOM.** You're serving a website. Your markdown is now HTML. You didn't need to configure webpack, you didn't need to set up a CI/CD pipeline, you didn't need to attend a 3-hour meeting about component architecture. You just ran one fucking command.

### "But what about my custom CSS?"

Oh, you want custom CSS? Drop a `style.css` file in `assets/style.css` and **it just fucking works.** No sass compilation, no postcss plugins, no "CSS-in-JS-in-CSS-in-JSON" framework horseshit. Just regular ass CSS that you learned in 2005 and still works perfectly fine.

```
your-site/
└── assets/
    └── style.css  ← Put your CSS here to override the default
```

That's the **only** special thing about the `assets/` folder - if `assets/style.css` exists, it replaces the built-in CSS. Everything else? You can organize however the fuck you want. Put images in `/images`, put PDFs in `/documents`, put memes in `/shitposts` - the server doesn't give a shit.

If you don't have `assets/style.css`, mofo uses a built-in minimal stylesheet that **doesn't look like shit** and **weighs fuck-all.** Unlike your favorite React component library that's somehow 400KB before you even render "Hello World."

### The motherfucking philosophy

This tool was built with the **[motherfucking website](https://motherfuckingwebsite.com/)** philosophy:

- **Fast** → No JavaScript unless YOU add it (why the fuck would you?)
- **Simple** → It's markdown. You write text. It becomes HTML. Congratulations, you're a web developer now.
- **Minimal** → The whole thing is one Go binary. No node_modules directory that achieves sentience.
- **Accessible** → Semantic HTML, readable text, actual fucking content instead of div soup.

Perfect for blogs, documentation, personal sites, or literally anything where you want to **share information** instead of showing off your ability to configure webpack.

### Features (that actually matter)

✅ **Serves markdown as HTML** - shocking, I know  
✅ **Serves everything else as static files** - images, CSS, PDFs, fonts, whatever the fuck you want  
✅ **Frontmatter support** - YAML/TOML for titles and meta tags  
✅ **GitHub Flavored Markdown** - tables, strikethrough, task lists, the good shit  
✅ **Zero configuration** - organize files however you want, there's no "right" structure  
✅ **Auto-heading IDs** - for anchor links, because we're not savages  
✅ **Hot-reload friendly** - crashes and retries, like your last relationship  

### Installation

```bash
go install github.com/NicoNex/mofo@latest
```

Or build it yourself if you don't trust binaries from strangers (smart):

```bash
git clone https://github.com/NicoNex/mofo.git
cd mofo
go build
```

Now you have a single binary. Move it wherever the fuck you want. `/usr/local/bin`, your home directory, a USB stick, I don't care. It doesn't need a `node_modules` folder the size of the goddamn moon.

### Usage

#### Serve a directory

```bash
mofo /path/to/site
```

Your site needs an `index.md` file. That's it. That's the requirement. Can you create one file? Great, you're qualified.

#### Serve on a different port

```bash
mofo -p 3000 /path/to/site
```

Or use `-port` if you're feeling verbose. I included both flags because I'm not a monster.

**IMPORTANT:** Flags go BEFORE the path, not after. This is how command-line tools work. If you put the path first, the flags get ignored, and then you'll bitch about the port not working. Don't be that person.

```bash
# ✅ Correct (you're not an idiot)
mofo -p 3000 /path/to/site

# ❌ Wrong (you're about to file a bug report that I'll close)
mofo /path/to/site -p 3000
```

### Directory Structure (or lack thereof)

Here's the deal: **organize your shit however you want.** There's literally ONE requirement: an `index.md` file. That's it. Everything else is up to you.

```
your-site/
├── index.md              ← Required. That's the list.
├── about.md              ← Put markdown files wherever
├── blog/
│   └── post.md           ← Subdirectories? Sure.
├── images/
│   └── cat.jpg           ← Static files? Anywhere you want.
├── downloads/
│   └── resume.pdf        ← PDFs, fonts, whatever
└── assets/               ← Special folder (explained below)
    └── style.css         ← ONLY special for overriding default CSS
```

### Here's how file serving actually works

**Markdown files (`.md`)** → Rendered as HTML with your template
**Everything else** → Served directly as static files

That's it. That's the whole algorithm. You want an image at `/images/cat.jpg`? Put `cat.jpg` in an `images/` folder. You want a PDF at `/docs/resume.pdf`? Put it in a `docs/` folder. There's no magic, no configuration, no "approved" directory structure. It's a file server. It serves files.

### The `assets/` folder (the only special thing)

The ONLY special behavior is `assets/style.css` - if this file exists, it **overrides the built-in CSS**. That's the entire special case. Everything else in `assets/` (or any other folder) is just... a file. Want images in `assets/`? Sure. Want them in `images/`? Also fine. Want them in `pictures-of-my-cat/`? I'm not your dad.

**Examples:**
- `![Cat](/images/cat.jpg)` → put cat.jpg in `images/`
- `![Dog](/assets/dog.png)` → put dog.png in `assets/`
- `<script src="/js/app.js">` → put app.js in `js/`
- `[Resume](/resume.pdf)` → put resume.pdf in the root

**It just serves the fucking files.** No preprocessing, no optimization, no "asset pipeline." If the file exists at that path, it gets served. If it doesn't, you get a 404. Revolutionary simplicity.

### Frontmatter (because metadata is occasionally useful)

```markdown
---
title: My Thoughts on JavaScript Frameworks
description: Spoiler alert - we have too many
---

# They keep making new ones

Why? Nobody knows...
```

Supported formats: **YAML** and **TOML**

Supported fields:
- `title` - Page title (shown in browser tab and `<title>` tag)
- `description` - Meta description (rendered in `<meta name="description">` tag for SEO)

If you don't provide a title, you get "Untitled" - which is perfect for your half-finished blog posts that you'll never publish. If you don't provide a description, no meta description tag is rendered.

**Favicon:** Just put `favicon.ico` in your site's root directory. That's it. No configuration needed. The template looks for `/favicon.ico` automatically. If the file exists, browsers show it. If it doesn't exist, browsers silently fail and you get no icon. Simple.

### How It Works (the technical shit)

1. Starts an HTTP server (port `:8080` by default)
2. Request comes in for a file
3. **Is it a `.md` file?** → Convert to HTML using [goldmark](https://github.com/yuin/goldmark), render with template
4. **Is it anything else?** → Serve it directly as a static file (images, CSS, JS, PDFs, whatever)
5. **File doesn't exist?** → 404, obviously
6. Auto-retries if it crashes (every 5 seconds, like a persistent ex)

**No build step.** No caching layer. No CDN. No GraphQL API to query your markdown files like it's 2019 and you're trying to get VC funding. Just a server that checks file extensions and serves files. Revolutionary, I know.

### Examples

#### Make a personal site in 30 seconds

```bash
mkdir my-site && cd my-site
echo -e "---\ntitle: Home\n---\n\n# Hi\n\nThis is my website.\n\nIt loads in 0.3 seconds." > index.md
mofo .
```

Open `http://localhost:8080` - you now have a website that's faster than 95% of the internet. You're welcome.

#### Add custom styling

```bash
mkdir -p assets
cat > assets/style.css << 'EOF'
body {
    max-width: 650px;
    margin: 40px auto;
    font-size: 18px;
    line-height: 1.6;
    color: #333;
}
EOF
```

Congratulations, you just wrote CSS that will still work in 2040, unlike whatever JavaScript framework you used last week.

### FAQ

**Q: Can I use this for production?**  
A: If your "production" is a blog, documentation site, or personal page - fuck yes. If you're building the next Facebook, probably use something else (but also, please don't build the next Facebook).

**Q: What about SEO?**  
A: You have semantic HTML, proper title tags, meta descriptions, and fast load times. That's literally all SEO is. Everything else is snake oil sold by people with "growth hacker" in their LinkedIn bio.

**Q: Does it support [insert framework here]?**  
A: No. And that's the whole fucking point.

**Q: What about dynamic content?**  
A: Write a different tool. This one serves markdown files. If you need a database and authentication and server-side rendering and GraphQL and microservices, you're in the wrong place, friend.

**Q: The default CSS is ugly.**  
A: Then change it. Put your CSS in `assets/style.css` and make it as ugly as you want. Or as pretty. I'm not your dad.

**Q: Can I add JavaScript?**  
A: You can, but *should you?* (No. The answer is no.)

### Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Actually good markdown parser
- [goldmark/frontmatter](https://go.abhg.dev/goldmark/frontmatter) - Frontmatter extension

That's it. Two dependencies. Both written in Go. Both compile to native code. No supply chain attacks from `leftpad-v2-react-hooks` maintained by a teenager in Belarus.

### Contributing

Keep it simple. Keep it fast. Don't add features "just in case." If you want to add React support, I will personally close your PR and suggest therapy.

### License

GPL3 - You can use it, modify it, sell it, whatever. But if you distribute your modified version, you gotta share the source code too. It's called copyleft, and it keeps shit free and open. Don't like it? Write your own markdown server.

---

**Built for people who miss the internet before it needed 3GB of RAM to display a recipe.**

*No JavaScript was harmed in the making of this tool.*
