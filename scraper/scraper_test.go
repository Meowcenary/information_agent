package scraper

import (
	"reflect"
	"testing"
)

func TestReadWikiJson(t *testing.T) {
	expectedPages := []WikiPage{
		WikiPage{
			Url: "testing.com",
			Title: "Test Title",
			Text: "Example text that might be stored into a WikiPage struct. Not particularly interesting",
		},
		WikiPage{
			Url: "testing_again.com",
			Title: "Another~Test~Title",
			Text: "Robotics is a very interesting field of study. The articles that are scraped in main have shown me this",
		},
		WikiPage{
			Url: "testing_the_third.com",
			Title: "Third Test Title!!",
			Text: "Imagine having a third example of text to store into a WikiPage struct",
		},
	}
	filepath := "../testdata/test_wikipage.json"
	pages, err := ReadWikiJson(filepath)

	if err != nil {
		t.Errorf("Error raised by ReadWikiJson: %s", err)
	} else if reflect.DeepEqual(pages, expectedPages) == false {
		t.Errorf("Values do not match what was expected\nWanted: %s\nGot: %s", expectedPages, pages)
	}
}

func TestWriteWikiJson(t *testing.T) {
	wikijson := []WikiPage{
		WikiPage{
			Url: "testjsonwriting.com",
			Title: "Test JSON Writing ",
			Text: "Testing out writing basic JSON to a file. Not the most interesting stuff.",
		},
		WikiPage{
			Url: "evenmoretesting.org",
			Title: "Even More Testing",
			Text: "You really can't have enough test data. Even stuff as simple as this is helpful to ensure things are working.",
		},
		WikiPage{
			Url: "lastwritingtestelement.net",
			Title: "Last Element",
			Text: "The last json writing test element. To make th\nings interes7ing there are s0me *w3ird* (h^ract3r$ in her3+!-",
		},
	}
	filepath := "../testdata/test_wikipage_write.json"
	err := WriteWikiJson(filepath, wikijson, false)

	if err != nil {
		t.Errorf("Error raised by WriteWikiJson: %s", err)
	}

	// Read data back from file
	pages, err := ReadWikiJson(filepath)

	if err != nil {
		t.Errorf("Error raised by ReadWikiJson while trying to read file written by WriteWikiJson: %s", err)
	} else if reflect.DeepEqual(pages, wikijson) == false {
		t.Errorf("Data read back from written file did not match what was expected\nWanted: %s\nGot: %s", wikijson, pages)
	}
}
