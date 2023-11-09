package main

import (
	"net/http"
	"time"
)

func NewNowHandler(now func() time.Time) NowHandler {
	return NowHandler{Now: now}
}

type NowHandler struct {
	Now func() time.Time
}

func (nh NowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timeComponent(nh.Now()).Render(r.Context(), w)
}

type DerpHandler struct {}

func (dh DerpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	derpComponent().Render(r.Context(), w)
}

func main() {
	http.Handle("/", NewNowHandler(time.Now))
	http.Handle("/derp", DerpHandler{})

	http.ListenAndServe(":8080", nil)
}
