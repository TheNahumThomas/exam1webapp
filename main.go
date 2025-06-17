package main

import (
	utilities "exam1webapp/utils"
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

	mux.Handle("/login", safehttp.MethodPost, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\login.html"))

		formData, err := r.PostForm()
		if err != nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		username := formData.String("username", "")
		password := formData.String("password", "")

		err = utilities.Login(username, password, w)
		if err != nil {
			return safehttp.ExecuteTemplate(w, t, map[string]interface{}{"Error": "Invalid credentials"})
		}

		data := map[string]string{"Username": username}

		return safehttp.ExecuteTemplate(w, t, data)
	}))

	mux.Handle("/upload", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html"))
		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	http.ListenAndServe(":8080", mux)
}
