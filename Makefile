GOPATH := $(shell go env GOPATH)

default: build

build: server.go components_templ.go
	go build -o server server.go components_templ.go

clean:
	-rm server

clean_wiki_pages:
	-rm -f wiki_page_json/*
	# recreate .gitignore
	echo "*" > wiki_page_json/.gitignore

scrape_wiki_pages:
	go build run_scraper.go
	./run_scraper

run:
	$(GOPATH)/bin/templ generate
	-rm server
	go build -o server server.go components_templ.go
	./server
