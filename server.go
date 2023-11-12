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
	http.Handle("/home", NewHomeHandler())
	http.Handle("/wiki_page_json/", NewPagesHandler())
	http.Handle("/search", NewSearchHandler())
	http.Handle("/about", NewAboutHandler())
	http.Handle("/scrape_wikipedia", NewScrapeWikipediaHandler())

	// Start the server.
	log.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

// Data Getters

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
	hh.Log.Println("getting pages from wiki_page_json/")
	pages, err := getPages("wiki_page_json")

	if err != nil {
		hh.Log.Println("failed to get pages: %v", err)
		return
	}

	hh.Log.Println("rendering home")
	templ.Handler(home(pages)).ServeHTTP(w, r)
}

// Scrape Wikipedia Handler
func NewScrapeWikipediaHandler() ScrapeWikipediaHandler {
	return ScrapeWikipediaHandler {
		Log: log.Default(),
	}
}

type ScrapeWikipediaHandler struct {
	Log *log.Logger
}

func (swh ScrapeWikipediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	swh.Log.Println("scraping url: ", url)
	page := scraper.ScrapeWikiUrls([]string{url})[0]
	filepath := "wiki_page_json" + "/" + page.FilenameFromTitle()
	log.Println("Writing to", filepath)
	scraper.WriteWikiPageJson(filepath, page)
	pages, _ := getPages("wiki_page_json")
	// Render Home
	templ.Handler(home(pages)).ServeHTTP(w, r)
}

// Search Handler

func NewSearchHandler() SearchHandler {
	return SearchHandler {
		Log: log.Default(),
	}
}

type SearchHandler struct {
	Log *log.Logger
}

func SearchPostHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Default()
	logger.Println("search post handler")
	var queryResults []scraper.WikiQueryResult
	r.ParseForm()

	var query string
	if r.Form.Has("search") {
		query = r.Form["search"][0]
		logger.Println("searching wikipedia with query: ", query)
		queryResults = scraper.SearchWikipedia(query)
	}

	logger.Println("rendering search results")
	templ.Handler(searchResults(query, queryResults)).ServeHTTP(w, r)
}

func (sh SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		sh.Log.Println("running passthrough search")
		SearchPostHandler(w, r)
		return
	}

	sh.Log.Println("rendering search")
	templ.Handler(search()).ServeHTTP(w, r)
}

// About Handler

func NewAboutHandler() AboutHandler {
	return AboutHandler {
		Log: log.Default(),
	}
}

type AboutHandler struct {
	Log *log.Logger
}

func (ah AboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.Log.Println("rendering about")
	templ.Handler(about()).ServeHTTP(w, r)
}

// Page Handler

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
	log.Println("retrieving page from filepath: ", filepath)
	wikipage, err := getPage(filepath)

	if err != nil {
		ph.Log.Printf("failed to get page: %v", err)
		http.Error(w, "failed to retrieve pages", http.StatusInternalServerError)
		return
	}

	log.Println("formatting page html")
	html := FormatPageHtml(wikipage)
	html = RemoveScriptTags(html)
	html = RemoveAnchorTags(html)

	log.Println("creating templ component")
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
