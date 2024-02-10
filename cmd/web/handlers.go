package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dinizgab/snippetbox/internal/models"
)

func (a Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		a.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	title := ":D titulo :D"
	content := "12341234"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (a Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id < 1 || err != nil {
		a.clientError(w, http.StatusBadRequest)

		return
	}
	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}

		return
	}

	data := a.newTemplateData(r)
    data.Snippet = snippet

    a.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (a Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)

		return
	}

	latestSnippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
        
        return
	}

    data := a.newTemplateData(r)
    data.Snippets = latestSnippets

    a.render(w, http.StatusOK, "home.tmpl.html", data)
}
