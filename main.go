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
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("root.html")
    t.Execute(w, nil)
}
