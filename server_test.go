package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Meowcenary/information_agent/scraper"
)

func TestHomeHandlerServeHTTP(t *testing.T) {
	// Create a new instance of AboutHandler
	hh := NewAboutHandler()
	// Create a fake HTTP request
	req, err := http.NewRequest("GET", "/home", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a fake HTTP response recorder
	w := httptest.NewRecorder()

	// Call the ServeHTTP method to handle the request
	hh.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, w.Code)
	}
}

func TestAboutHandlerServeHTTP(t *testing.T) {
	// Create a new instance of AboutHandler
	ah := NewAboutHandler()
	// Create a fake HTTP request
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a fake HTTP response recorder
	w := httptest.NewRecorder()

	// Call the ServeHTTP method to handle the request
	ah.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, w.Code)
	}
}

func TestFormatPageHtml(t *testing.T) {
	// Create a sample WikiPage
	wikiPage := scraper.WikiPage{
		Title: "Test Title",
		Paragraphs: []scraper.WikiPageParagraph{
			{Text: "Test Paragraph 1"},
			{Text: "Test Paragraph 2"},
		},
	}

	// Test the FormatPageHtml function
	result := FormatPageHtml(wikiPage)

	// Expected HTML for the sample WikiPage
	expectedHTML := "<html><body><h1>Test Title</h1><hr></hr><p>Test Paragraph 1</p><p>Test Paragraph 2</p></body></html>"

	// Check that the result matches the expected value
	if result != expectedHTML {
		t.Errorf("Expected HTML:\n%s\n\nGot HTML:\n%s", expectedHTML, result)
	}
}

func TestRemoveScriptTags(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"<script>alert('Hello, World!');</script>", ""},
		{"<p>This is some text.</p><script>console.log('Script');</script><div>More text</div>", "<p>This is some text.</p><div>More text</div>"},
		{"<script src='external.js'></script>", ""},
		{"<div><script>console.log('Nested Script');</script></div>", "<div></div>"},
	}

	for _, tc := range testCases {
		result := RemoveScriptTags(tc.input)

		// Check that the result matches the expected value
		if result != tc.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}

func TestRemoveAnchorTags(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"<a href='https://example.com'>Link Text</a>", "Link Text"},
		{"<p>This is <a href='#'>a link</a>.</p>", "<p>This is a link.</p>"},
		{"<a href='https://example.org'>Another Link</a> <span>Text</span>", "Another Link <span>Text</span>"},
		{"<a href='#'>Click me!</a><div><a href='#'>Another Link</a></div>", "Click me!<div>Another Link</div>"},
	}

	for _, tc := range testCases {
		result := RemoveAnchorTags(tc.input)

		// Check that the result matches the expected value
		if result != tc.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}
