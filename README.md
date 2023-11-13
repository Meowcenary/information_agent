# Applied Statistics Method Reference
Prototype information agent web application written in Go. This specific instance is configured with a focus on applied
statistics.

### Installing the App
The app can be installed with git: `git clone git@github.com:Meowcenary/information_agent.git` or
`git clone https://github.com/Meowcenary/information_agent.git`.

### Running the Web Scraper
From the root directory: `make scrape_wiki_pages` will scrape the URLs specified in the new line delimited file
`urls.txt`. The default URLs are set with a focus on data science topics, but this can be edited to scrape other pages
on Wikipedia.

### Running the Web App
From the root directory: `make run`. Once the server has started navigate to `http:/localhost:8000/home`. The app
supports viewing wiki pages, searching wikipedia, and local file CRUD operations. It is recommended to run the scraper
before running the web app to seed the content.

### Demo Workflow
From the root directory run:
```
make scrape_wiki_pages
make run
```

Once the server is running navigate to http:/localhost:8000/home to see the list of scraped pages. Click "View Topic" to
view the scraped data formatted for local viewing. Click "Delete Topic" to remove a page from the system. On the top nav
bar click "Search" and run a simple search for: "query tests" to see a list of articles that could be scraped from Wikipedia.
The button "View Article on Wikipedia" will open a new tab to the article on Wikipedia and the button "Add to Reference"
will scrape the page and add the link to the home page list.

### Running the Tests
To run the full test suite run the following from the root directory: `make test`

### Updating Views
This project uses the Go package "templ" to build it's views. Most of the view code is in `components.templ`, but it is
also possible to build templ components with Go code formatting the html. It should be noted that this is not very
secure since it opens up the potential for cross site scripting, but in this instance the html returned from the scraper
is from a trusted source.

### Third Party Go Packages Used
- [Templ](https://pkg.go.dev/github.com/josegpt/go-utils/templ) - template package for Go. Documentation can be read [here](https://templ.guide/)
- [Colly](https://github.com/gocolly/colly) - scraper and crawler package for Go. Documentation can be read [here](https://go-colly.org/)
