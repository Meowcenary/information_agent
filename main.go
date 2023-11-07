package main

import (
  // "fmt"
  "html/template"
  "log"
  "net/http"
)

// type DataSource struct {
//   BaseUrl string // E.g en.wikipedia.org/wiki/
//   Params []string // E.g "Mediterranean_Sea"
// }

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/submit", submitHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

// handler functions must have http.RepsonseWriter and *http.Request as function params
func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("root.html")
    t.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
  // handle the query here
  t, _ := template.ParseFiles("query_results.html")
  t.Execute(w, nil)
}
