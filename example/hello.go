package main

import (
	"github.com/pat-go/pat2"
	"io"
	"net/http"
)

func main() {
	m := pat.New()
	m.Get("/hello/:name", pat.HandlerFunc(hello))
	m.Get("/splat/", pat.HandlerFlat(splat))
	http.ListenAndServe("localhost:5000", m)
}

func hello(params pat.Params, _ string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := params[":name"]
		io.WriteString(w, "Hello, "+name)
	})
}

func splat(_ pat.Params, s string, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Splat: "+s)
}
