package main

import "net/http"

func (a Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippets/view", a.snippetView)
	mux.HandleFunc("/snippets/create", a.snippetCreate)

	return mux
}
