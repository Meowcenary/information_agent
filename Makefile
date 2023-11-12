default: build

build: server.go components_templ.go
	go build -o server server.go components_templ.go

clean:
	-rm server

run:
	~/go/bin/templ generate
	-rm server
	go build -o server server.go components_templ.go
	./server
