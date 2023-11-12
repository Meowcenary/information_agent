package main

import (
	"log"

	"github.com/Meowcenary/information_agent/scraper"
)

// Scrape wikipedia urls and store result into
// json files to be used by the web app
func main() {
	urls := []string{
		"https://en.wikipedia.org/wiki/Bootstrapping_(statistics)",
		"https://en.wikipedia.org/wiki/Cross-validation_(statistics)",
		"https://en.wikipedia.org/wiki/Sampling_(statistics)",
		"https://en.wikipedia.org/wiki/Resampling_(statistics)",
	}
	dataDir := "wiki_page_json"
	pages := scraper.ScrapeWikiUrls(urls)
	log.Printf("Done scraping, writing to files")

	for _, page := range pages {
		filepath := dataDir + "/" + page.FilenameFromTitle()
		log.Println("Writing to", filepath)
		scraper.WriteWikiPageJson(filepath, page)
	}
}
