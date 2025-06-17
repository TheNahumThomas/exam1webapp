package utilities

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"golang.org/x/crypto/bcrypt"
)

const dbFile = "webapp.db"

const loginQuery = "SELECT `password` FROM tbl_users WHERE `username` = ?"
const sessionStoreQuery = "UPDATE `tbl_users` SET `sessionid` = ? WHERE `username` = ?"
const checkSessionIDQuery = "SELECT `username` FROM `tbl_users` WHERE `sessionid` = ?"

func dbConnect() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func Login(username, password string) (string, error) {
	var storedPassword string

	db, err := dbConnect()
	if err != nil {
		return "", err
	}

	defer db.Close()

	err = db.QueryRow(loginQuery, username).Scan(&storedPassword)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(sessionStoreQuery, sessionID, username)
	if err != nil {
		return "", err
	}
	fmt.Println(sessionID)
	return sessionID, nil
}

func CheckSession(sessionID string) (string, error) {
	db, err := dbConnect()
	if err != nil {
		return "", err
	}

	var username string
	err = db.QueryRow(checkSessionIDQuery, sessionID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return username, nil
}
