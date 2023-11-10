package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	log.Println("starting...")
	// Http handlers
	http.Handle("/", NewHomeHandler())
	http.Handle("/topics", NewTopicsHandler())

	// Start the server.
	log.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

func getTopics() ([]Topic, error) {
	topics := []Topic{
		Topic{Name: "Test Topic", Body: "test body"},
		Topic{Name: "Test Topic 2", Body: "another test body"},
	}

	return topics, nil
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
	topics, err := getTopics()

	if err != nil {
		hh.Log.Printf("failed to get posts: %v", err)
		http.Error(w, "failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	hh.Log.Println("rendering home")
	templ.Handler(home(topics)).ServeHTTP(w, r)
}

// Topics Handler
func NewTopicsHandler() TopicsHandler {
	return TopicsHandler {
		GetTopics: getTopics,
		Log: log.Default(),
	}
}

type TopicsHandler struct {
	GetTopics func() ([]Topic, error)
	Log *log.Logger
}

func (th TopicsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	topics, err := getTopics()

	if err != nil {
		th.Log.Printf("failed to get topics: %v", err)
		http.Error(w, "failed to retrieve topics", http.StatusInternalServerError)
		return
	}

	th.Log.Println("rendering home")
	templ.Handler(home(topics)).ServeHTTP(w, r)
}

type Topic struct {
	Name string
	Body string
}
