# Applied Statistics Method Reference
Prototype information agent web application written in Go. This specific instance is configured with a focus on applied
statistics.

### Running the Web Scraper
From the root directory run: `go run run_scraper.go`. To update the page that are scraped edit the new line delmited
file `urls.txt`.

### Running the Web App
From the root directory run: `go run main.go component_templ.go`. Once the server is running open a web browser and
navigate to "http://localhost:8000".

### Updating Views
This project uses the Go package "templ" to build it's views. Most of the view code is in `components.templ`, but it is
also possible to build templ components with Go code formatting the html. It should be noted that this is not very
secure since it opens up the potential for cross site scripting, but in this instance the html returned from the scraper
is from a trusted source.

### Third Party Go Packages Used
- [Templ](https://pkg.go.dev/github.com/josegpt/go-utils/templ) - template package for Go. Documentation can be read [here](https://templ.guide/)
- [Colly](https://github.com/gocolly/colly) - scraper and crawler package for Go. Documentation can be read [here](https://go-colly.org/)
