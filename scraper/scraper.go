package scraper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

// Struct to hold scraped data
type WikiPage struct {
	Url string `json:"url"`
	Title string `json:"title"`
	Text string `json:"text"`
	Tags []string `json:"tags"`
}

// Constructor
func NewWikiPage(url string) *WikiPage {
	return &WikiPage{
		Url: url,
		Title: "",
		Text: "",
		Tags: []string{},
	}
}

// Read json created from scraper
func ReadWikiJson(filepath string) ([]WikiPage, error) {
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()

	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(jsonFile)
	var pages []WikiPage
	err = dec.Decode(&pages)

	return pages, err
}

// In general single line is preferred to save memory, but for
// the assignment newline delimited json is set to default
func WriteWikiJson(filepath string, wikijson []WikiPage, newlineDelim bool) error {
	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)

	if newlineDelim == true {
		file.WriteString("[")

		for _, record := range wikijson {
			enc.Encode(record)

			if err != nil {
				break
			}
		}

		// Remove last newline char and write ] char
		file.Seek(-1, 1)
		file.WriteString("]")
	} else {
		err = enc.Encode(wikijson)
	}

	return err
}

func ScrapeWikiUrls(urls []string) []WikiPage {
	// Pointer to WikiPage object, used to work around passing arguments to callbacks
	var currentPage *WikiPage
	// Json of scraped wikipage data to return
	var json []WikiPage
	// Colly
	wikiPageCollector := colly.NewCollector()
	wikiPageCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	wikiPageCollector.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
		currentPage.Title = fmt.Sprintf("%s", e.Text)
	})
	// "#mw-content-text" is taken from Python script "//div[@id="mw-content-text"]"
	wikiPageCollector.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		currentPage.Text = fmt.Sprintf("%s", e.Text)
	})

	for _, url := range urls {
		// Assign currentPage pointer to a new WikiPage struct
		currentPage = NewWikiPage(url)
		// Visit url and trigger scraping functions
		wikiPageCollector.Visit(url)
		// Dereference and append WikiPage struct
		json = append(json, *currentPage)
	}

	return json
}
