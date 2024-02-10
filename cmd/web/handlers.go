package main

import (
	"errors"
	"fmt"
	"html/template"
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

	files := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	snippetData := &templateData{Snippet: snippet}

	err = ts.ExecuteTemplate(w, "base", snippetData)
	if err != nil {
		a.serverError(w, err)
	}

	fmt.Fprintf(w, "snippet: %+v", snippet)
}

func (a Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)

		return
	}

	latestSnippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
	}

	for _, snip := range latestSnippets {
		fmt.Fprintf(w, "%+v\n", snip)
	}

	// files := []string{
	// 	"./ui/html/pages/home.tmpl.html",
	// 	"./ui/html/pages/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	a.errorLog.Println(err.Error())
	// 	a.serverError(w, err)

	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	a.errorLog.Println(err.Error())
	// 	a.serverError(w, err)
	// }

	w.Write([]byte("This is the home page!!!"))
}
