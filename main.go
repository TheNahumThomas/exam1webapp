package main

import (
	utilities "exam1webapp/utils"
	"fmt"
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
		authenticated := false
		username := "Unknown User"
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\index.html"))
		userSessionID, err := r.Cookie("session")
		if err != nil {
			fmt.Println("Error retrieving cookies:", err)
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			fmt.Println("Error retrieving cookies:", err)
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		if userSessionID == nil || usernameCookie == nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		username = usernameCookie.Value()
		sessionOnwer, err := utilities.CheckSession(userSessionID.Value())
		if err != nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		if sessionOnwer != username {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		} else {
			authenticated = true
		}

		if !authenticated {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		return safehttp.ExecuteTemplate(w, t, map[string]interface{}{"Username": username})
	}))

	mux.Handle("/login", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\login.html"))
		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	mux.Handle("/login/attempt", safehttp.MethodPost, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\login.html"))
		formData, err := r.PostForm()
		if err != nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		username := formData.String("username", "")
		password := formData.String("password", "")

		sessionID, err := utilities.Login(username, password)
		fmt.Println(err)
		if err != nil {
			return safehttp.ExecuteTemplate(w, t, map[string]interface{}{"Error": "Invalid username or password"})
		}

		sessionCookie := safehttp.NewCookie("session", sessionID)
		usernameCookie := safehttp.NewCookie("username", username)
		usernameCookie.Path("/")
		sessionCookie.Path("/")

		err = safehttp.ResponseHeadersWriter(w).AddCookie(sessionCookie)
		if err != nil {
			return safehttp.ExecuteTemplate(w, t, map[string]interface{}{"Error": "Failed to set session cookies"})
		}
		err = safehttp.ResponseHeadersWriter(w).AddCookie(usernameCookie)
		if err != nil {
			return safehttp.ExecuteTemplate(w, t, map[string]interface{}{"Error": "Failed to set session cookies"})
		}

		return safehttp.Redirect(w, r, "/home", safehttp.StatusFound)

	}))

	mux.Handle("/patient-record", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		t := template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html"))
		userSessionID, err := r.Cookie("session")
		if err != nil {
			fmt.Println("Error retrieving cookies:", err)
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			fmt.Println("Error retrieving cookies:", err)
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		if userSessionID == nil || usernameCookie == nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		username := usernameCookie.Value()
		sessionOnwer, err := utilities.CheckSession(userSessionID.Value())
		if err != nil {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		if sessionOnwer != username {
			return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
		}

		return safehttp.ExecuteTemplate(w, t, nil)
	}))

	mux.Handle("/logout", safehttp.MethodGet, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		sessionReset := safehttp.NewCookie("session", "")
		sessionReset.Path("/")
		err := safehttp.ResponseHeadersWriter(w).AddCookie(sessionReset)
		if err != nil {
			return safehttp.ExecuteTemplate(w, template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\login.html")), map[string]interface{}{"Error": "Failed to logout"})
		}
		return safehttp.Redirect(w, r, "/login", safehttp.StatusFound)
	}))

	mux.Handle("/upload", safehttp.MethodPost, safehttp.HandlerFunc(func(w safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		formData, err := r.MultipartForm(20000)
		if err != nil {
			return safehttp.ExecuteTemplate(w, template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html")), map[string]interface{}{"Error": "Failed to upload file"})
		}

		file := formData.File("file")
		if len(file) == 0 {
			return safehttp.ExecuteTemplate(w, template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html")), map[string]interface{}{"Error": "No file uploaded"})
		}

		file[0].Open()

		return safehttp.ExecuteTemplate(w, template.Must(template.ParseFiles("static\\templates\\base.html", "static\\templates\\upload.html")), map[string]interface{}{"Success": "File uploaded successfully"})

	}))

	http.ListenAndServe(":8080", mux)
}
