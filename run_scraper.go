package main

import (
	"log"

	"github.com/Meowcenary/information_agent/scraper"
)

// Scrape wikipedia urls and store result into
// json files to be used by the web app
func main() {
	urls, _ := scraper.ReadUrlsFromTextFile("urls.txt")
	dataDir := "wiki_page_json"
	pages := scraper.ScrapeWikiUrls(urls)
	log.Printf("Done scraping, writing to files")

	for _, page := range pages {
		filepath := dataDir + "/" + page.FilenameFromTitle()
		log.Println("Writing to", filepath)
		scraper.WriteWikiPageJson(filepath, page)
	}
}
