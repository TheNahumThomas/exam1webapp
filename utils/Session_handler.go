package utilities

import (
	"net/http"

	"github.com/google/go-safeweb/safehttp"
)

type SessionData struct {
	Username    string
	Permissions []string
}

func AttemptLogin(w safehttp.ResponseWriter, r *safehttp.IncomingRequest, username, password string) safehttp.Result {

	err := login(username, password, w)
	if err != nil {
		return safehttp.Error(http.StatusUnauthorized, "Invalid credentials")
	}

	sessionData := &SessionData{
		Username:    username,
		Permissions: []string{"read", "write"}, // Example permissions
	}
	session.Set("sessionData", sessionData)

	return safehttp.Redirect(w, r, "/home", safehttp.StatusFound)
}
