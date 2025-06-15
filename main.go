package main

import (
	"net/http"

	"github.com/google/safehtml/template"

	"github.com/google/go-safeweb/safehttp"
)

func main() {
	mux := safehttp.NewServeMuxConfig(safehttp.DefaultDispatcher{}).Mux()

	mux.Handle("/", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		return safehttp.Redirect(w, r, "/home", safehttp.StatusFound)
	}))

	mux.Handle("/home", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\index.html"))
		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	mux.Handle("/login", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\login.html"))
		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	mux.Handle("/upload", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html"))
		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	http.ListenAndServe(":8080", mux)
}
