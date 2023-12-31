package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

// WikiPageParagraph

type WikiPageParagraph struct {
	Text string `json:"text"`
}

// WikiPage

// Constructor
func NewWikiPage(url string) *WikiPage {
	return &WikiPage {
		Url: url,
		Title: "",
		Text: "",
		Tags: []string{},
	}
}

// Struct to hold scraped data
type WikiPage struct {
	Url string `json:"url"`
	Title string `json:"title"`
	Text string `json:"text"`
	Tags []string `json:"tags"`
	Paragraphs []WikiPageParagraph `json:"paragraphs"`
}

func ReadUrlsFromTextFile(filePath string) ([]string, error) {
	// Open the text file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffered reader.
	reader := bufio.NewReader(file)

	// Read the lines from the text file into a slice of strings.
	lines := []string{}
	for line, err := reader.ReadString('\n'); err != io.EOF; line, err = reader.ReadString('\n') {
		lines = append(lines, line[:len(line)-1])
	}

	return lines, nil
}

func FilenameFromTitle(title string) string {
	title = strings.ToLower(title)
	replacer := strings.NewReplacer(" ", "_", "(", "",  ")", "", "-", "")
	title = replacer.Replace(title)

	return title + ".json"
}

func (wp WikiPage) FilenameFromTitle() string {
	title := strings.ToLower(wp.Title)
	replacer := strings.NewReplacer(" ", "_", "(", "",  ")", "", "-", "")
	title = replacer.Replace(title)

	return title + ".json"
}

func DeleteWikiPageJson(filepath string) (error) {
	err := os.Remove(filepath)
	return err
}

// Read json created from scraper
// Returns pointer to WikiPage because using a WikiPage struct directly does not allow for a nil
// return. See here: https://stackoverflow.com/questions/50697914/return-nil-for-a-struct-in-go
func ReadWikiPageJson(filepath string) (*WikiPage, error) {
	var page WikiPage
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()

	if err != nil {
		// return page without data and an error
		return nil, err
	}

	dec := json.NewDecoder(jsonFile)
	err = dec.Decode(&page)

	return &page, err
}

// In general single line is preferred to save memory, but for
// the assignment newline delimited json is set to default
func WriteWikiPageJson(filepath string, page WikiPage) error {
	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	err = enc.Encode(page)

	return err
}

func ReadWikiPagesFromDirectory(dirPath string) ([]WikiPage, error) {
	// Get a list of all JSON files in the directory.
	files, err := ioutil.ReadDir(dirPath)

	if err != nil {
			return nil, err
	}

	// Parse all json pages
	wikiPages := make([]WikiPage, len(files))
	for index, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// open file and create decoder
		jsonFile, err := os.Open(filepath.Join(dirPath, file.Name()))
		dec := json.NewDecoder(jsonFile)

		if err != nil {
			return nil, err
		}

		// decode JSON into WikiPage struct
		var wikiPage WikiPage
		err = dec.Decode(&wikiPage)

		if err != nil {
			return nil, err
		}

		wikiPages[index] = wikiPage
	}

	return wikiPages, nil
}

func ScrapeWikiUrls(urls []string) []WikiPage {
	// Pointer to WikiPage object, used to work around passing arguments to callbacks
	var currentPage *WikiPage
	// Json of scraped wikipage data to return
	var pages []WikiPage

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
		e.ForEach("p", func(_ int, pe *colly.HTMLElement) {
			paragraphHtml, _ := pe.DOM.Html()
			paragraphText := fmt.Sprintf("%s", paragraphHtml)
			currentPage.Paragraphs = append(currentPage.Paragraphs, WikiPageParagraph{Text: paragraphText})
    })
	})

	for _, url := range urls {
		// Assign currentPage pointer to a new WikiPage struct
		currentPage = NewWikiPage(url)
		// Visit url and trigger scraping functions
		wikiPageCollector.Visit(url)
		// Dereference and append WikiPage struct
		pages = append(pages, *currentPage)
	}

	return pages
}

type WikiQueryResult struct {
	Url string `json:"url"`
	Title string `json:"title"`
}

// Search wikiepdia for query and return results as map of page title to page url
// Only works for one word query currently (lol)
func SearchWikipedia(query string) []WikiQueryResult {
	// Colly
	wikiSearchCollector := colly.NewCollector()
	wikiSearchCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	var queryResults []WikiQueryResult
	// If the query matches with a wikipedia page title it will go directly
	// to the page instead of listing the results of the search
	wikiSearchCollector.OnHTML(".mw-search-result-heading", func(e *colly.HTMLElement) {
		queryResult := WikiQueryResult{
			Url: "https://en.wikipedia.org/" + e.ChildAttr("a", "href"),
			Title: e.ChildAttr("a", "title"),
		}
		queryResults = append(queryResults, queryResult)
	})

	wikiSearchCollector.Visit("https://en.wikipedia.org/wiki/Special:Search?search=" + query)
	return queryResults
}


