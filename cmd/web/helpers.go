package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

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
