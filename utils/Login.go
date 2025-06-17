package utilities

import (
	"database/sql"
	"errors"

	"github.com/google/go-safeweb/safehttp"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"golang.org/x/crypto/bcrypt"
)

const dbFile = "exam1webapp/webapp.db"

const loginQuery = `SELECT password FROM tbl_users WHERE username = ?`

func dbConnect() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// func generateSessionID() (string, error) {
// 	b := make([]byte, 32)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(b), nil
// }

func Login(username, password string, w safehttp.ResponseWriter) error {
	var storedPassword string
	db, err := dbConnect()
	if err != nil {
		return err
	}

	err = db.QueryRow(loginQuery, username).Scan(&storedPassword)
	if err != nil {
		w.WriteError(safehttp.StatusUnauthorized)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}
