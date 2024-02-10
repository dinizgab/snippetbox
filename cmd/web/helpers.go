package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (a Application) newTemplateData(r *http.Request) *templateData {
    return &templateData{
        CurrentYear: time.Now().Year(),
    }
}

func (a Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a Application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
    ts, ok := a.templateCache[page]
    if !ok {
        err := fmt.Errorf("The template %s does not exists", page)
        a.serverError(w, err)

        return 
    }
    
    // First we load the file into a buffer, if there is any error we return the error
    // before loading the page in the screen
    buf := new(bytes.Buffer)
    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        a.serverError(w, err)

        return
    }

    w.WriteHeader(status)

    buf.WriteTo(w)
}
