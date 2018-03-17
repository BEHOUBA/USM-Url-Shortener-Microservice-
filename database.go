package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Initiation of the package by definition database connection setup
func init() {
	var err error
	db, err = sql.Open("postgres", "dbname=urlsdatabase user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// this methode store the urls pair struct into the database
func (urls *URLpair) storeUrls() (err error) {
	statement := "INSERT INTO urlmap (original, short_url) VALUES ($1, $2)"

	_, err = db.Exec(statement, urls.Original, urls.Short)

	if err != nil {
		fmt.Println(err)
	}
	return
}

// this function retrieve the corresponding long url from the database if it exist otherwise
// store the url pair to the database
func (urls *URLpair) getUrls() (err error) {
	statement := "SELECT ORIGINAL, SHORT_URL FROM URLMAP WHERE original=$1"

	err = db.QueryRow(statement, urls.Original).Scan(&urls.Original, &urls.Short)
	if err != nil {
		urls.storeUrls()
		return
	}
	fmt.Println(urls)
	return
}
