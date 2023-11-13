package scraper

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "testing"
)

func TestFilenameFromTitle(t *testing.T) {
	testCases := []struct {
		title    string
		expected string
	}{
		{"Go Programming Language", "go_programming_language.json"},
		{"Python (Programming Language)", "python_programming_language.json"},
		{"C++", "c++.json"},
		{"Title_With_Underscores", "title_with_underscores.json"},
	}

	for _, tc := range testCases {
		result := FilenameFromTitle(tc.title)

		// Check that the result matches the expected value
		if result != tc.expected {
			t.Errorf("For title '%s', expected filename '%s', got '%s'", tc.title, tc.expected, result)
		}
	}
}

func TestReadUrlsFromTextFile(t *testing.T) {
	// Create a temporary text file for testing
	tmpFile, err := ioutil.TempFile("", "testfile*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some sample URLs to the file
	urls := []string{
		"https://example.com",
		"https://example.org",
		"https://example.net",
	}
	for _, url := range urls {
		_, err := tmpFile.WriteString(url + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test the ReadUrlsFromTextFile function
	readUrls, err := ReadUrlsFromTextFile(tmpFile.Name())

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the number of read URLs matches the expected number
	if len(readUrls) != len(urls) {
		t.Errorf("Expected %d URLs, got %d URLs", len(urls), len(readUrls))
	}

	// Check that each read URL matches the expected URL
	for i, url := range urls {
		if readUrls[i] != url {
			t.Errorf("Expected URL %s, got %s", url, readUrls[i])
		}
	}
}

func TestDeleteWikiPageJson(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "testfile*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Test the DeleteWikiPageJson function
	err = DeleteWikiPageJson(tmpFile.Name())

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the file has been deleted
	_, err = os.Stat(tmpFile.Name())
	if !os.IsNotExist(err) {
		t.Errorf("Expected file to be deleted, but it still exists")
	}
}

func TestReadWikiPageJson(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "testfile*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a sample WikiPage to write to the file
	page := WikiPage{
		Title: "Test Title",
		Paragraphs: []WikiPageParagraph{
			{Text: "Test Paragraph 1"},
			{Text: "Test Paragraph 2"},
		},
	}

	// Marshal the WikiPage to JSON and write it to the file
	content, err := json.Marshal(page)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(tmpFile.Name(), content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test the ReadWikiPageJson function
	readPage, err := ReadWikiPageJson(tmpFile.Name())

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the read WikiPage matches the original
	if readPage.Title != page.Title {
		t.Errorf("Expected title %s, got %s", page.Title, readPage.Title)
	}

	// Check the number of paragraphs
	if len(readPage.Paragraphs) != len(page.Paragraphs) {
		t.Errorf("Expected %d paragraphs, got %d", len(page.Paragraphs), len(readPage.Paragraphs))
	}

	// Check the content of each paragraph
	for i, para := range page.Paragraphs {
		if readPage.Paragraphs[i].Text != para.Text {
			t.Errorf("Expected paragraph %s, got %s", para.Text, readPage.Paragraphs[i].Text)
		}
	}
}

func TestWriteWikiPageJson(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "testfile*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a sample WikiPage to write to the file
	page := WikiPage{
		Title: "Test Title",
		Paragraphs: []WikiPageParagraph{
			{Text: "Test Paragraph 1"},
			{Text: "Test Paragraph 2"},
		},
	}

	// Test the WriteWikiPageJson function
	err = WriteWikiPageJson(tmpFile.Name(), page)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Read the content of the written file
	content, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the JSON content back into a WikiPage
	var writtenPage WikiPage
	err = json.Unmarshal(content, &writtenPage)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the written WikiPage matches the original
	if writtenPage.Title != page.Title {
		t.Errorf("Expected title %s, got %s", page.Title, writtenPage.Title)
	}

	// Check the number of paragraphs
	if len(writtenPage.Paragraphs) != len(page.Paragraphs) {
		t.Errorf("Expected %d paragraphs, got %d", len(page.Paragraphs), len(writtenPage.Paragraphs))
	}

	// Check the content of each paragraph
	for i, para := range page.Paragraphs {
		if writtenPage.Paragraphs[i].Text != para.Text {
			t.Errorf("Expected paragraph %s, got %s", para.Text, writtenPage.Paragraphs[i].Text)
		}
	}
}

func TestReadWikiPagesFromDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary JSON file for testing
	tmpFile, err := ioutil.TempFile(tmpDir, "testfile*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some JSON content to the file
	jsonContent := `{"Title": "Test Title", "Paragraphs": [{"Text": "Test Paragraph"}]}`
	err = ioutil.WriteFile(tmpFile.Name(), []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test the ReadWikiPagesFromDirectory function
	wikiPages, err := ReadWikiPagesFromDirectory(tmpDir)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that at least one WikiPage was read
	if len(wikiPages) == 0 {
		t.Errorf("Expected at least one WikiPage, got none")
	}

	// Check that the title and paragraphs are correctly read
	expectedTitle := "Test Title"
	expectedParagraph := "Test Paragraph"
	if wikiPages[0].Title != expectedTitle {
		t.Errorf("Expected title %s, got %s", expectedTitle, wikiPages[0].Title)
	}
	if len(wikiPages[0].Paragraphs) == 0 || wikiPages[0].Paragraphs[0].Text != expectedParagraph {
		t.Errorf("Expected paragraph %s, got %s", expectedParagraph, wikiPages[0].Paragraphs[0].Text)
	}
}

func TestScrapeWikiUrls(t *testing.T) {
	urls := []string{
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
		"https://en.wikipedia.org/wiki/Python_(programming_language)",
	}

	pages := ScrapeWikiUrls(urls)

	// Check that the number of returned pages is equal to the number of input URLs
	if len(pages) != len(urls) {
		t.Errorf("Expected %d pages, got %d pages", len(urls), len(pages))
	}

	// Check that each returned page has a title and paragraphs
	for _, page := range pages {
		if page.Title == "" {
			t.Errorf("Page should have a non-empty title, got an empty title")
		}
		if len(page.Paragraphs) == 0 {
			t.Errorf("Page should have at least one paragraph, got none")
		}
	}
}

func TestSearchWikipedia(t *testing.T) {
	query := "stadfs"
	results := SearchWikipedia(query)

	// Check that there are results
	if len(results) == 0 {
		t.Errorf("Expected non-empty results for query '%s', got empty results", query)
	}

	// Check that each result has a title and URL
	for _, result := range results {
		if result.Title == "" {
			t.Errorf("Result should have a non-empty title, got an empty title")
		}
		if result.Url == "" {
			t.Errorf("Result should have a non-empty URL, got an empty URL")
		}
	}
}
