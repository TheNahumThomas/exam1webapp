package utilities

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Transaction struct {
	db *sql.DB
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Transaction) Login(username, password string, w http.ResponseWriter) error {
	var storedPassword string
	err := s.db.QueryRow("SELECT password FROM tbl_users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}
