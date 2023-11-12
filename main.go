package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	"github.com/Meowcenary/information_agent/scraper"
)

func main() {
	log.Println("starting...")

	// Http handlers
	http.Handle("/", NewHomeHandler())
	http.Handle("/wiki_page_json/", NewPagesHandler())

	// Start the server.
	log.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

// This is effectively an index of pages
func getPages(dirPath string) ([]scraper.WikiPage, error) {
	wikiPages, err := scraper.ReadWikiPagesFromDirectory(dirPath)

	return wikiPages, err
}

func getPage(filepath string) (scraper.WikiPage, error) {
	wikiPage, err := scraper.ReadWikiPageJson(filepath)

	return *wikiPage, err
}

// Home Handler
func NewHomeHandler() HomeHandler {
	return HomeHandler {
		Log: log.Default(),
	}
}

type HomeHandler struct {
	Log *log.Logger
}

func (hh HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages, err := getPages("wiki_page_json")

	if err != nil {
		hh.Log.Println("failed to get pages: %v", err)
		return
	}

	hh.Log.Println("rendering home")
	templ.Handler(home(pages)).ServeHTTP(w, r)
}

func NewPagesHandler() PagesHandler {
	return PagesHandler {
		Log: log.Default(),
	}
}

type PagesHandler struct {
	Log *log.Logger
}

func (ph PagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ph.Log.Println(path)
	filepath := scraper.FilenameFromTitle(r.URL.Path[1:])
	log.Println("Retrieving page from filepath: ", filepath)
	wikipage, err := getPage(filepath)

	if err != nil {
		ph.Log.Printf("failed to get page: %v", err)
		http.Error(w, "failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	// Create an html string to render
	html := FormatPageHtml(wikipage)
	html = RemoveScriptTags(html)
	html = RemoveAnchorTags(html)
	// Create an unsafe component containing raw HTML.
	content := Unsafe(html)
	ph.Log.Println("rendering pages")
	templ.Handler(page(content)).ServeHTTP(w, r)
}

// Page Component Building

func FormatPageHtml(wikipage scraper.WikiPage) string {
	html := "<html><body><h1>" + wikipage.Title + "</h1><hr></hr>"
	for _, paragraph := range wikipage.Paragraphs {
		html += "<p>" + paragraph.Text + "</p>"
	}
	html += "</body></html>"
	return html
}

// RemoveScriptTags removes script tags and their content from an HTML formatted string
func RemoveScriptTags(html string) string {
	scriptTagRegex := regexp.MustCompile(`<script(.*?)>(.*?)</script>`)
	html = scriptTagRegex.ReplaceAllString(html, "")
	return html
}

// RemoveAnchorTags removes anchor tags from an HTML formatted string but keeps the link text
func RemoveAnchorTags(html string) string {
	anchorTagRegex := regexp.MustCompile(`<a(.*?)>(.*?)</a>`)
	html = anchorTagRegex.ReplaceAllString(html, "$2")
	return html
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
