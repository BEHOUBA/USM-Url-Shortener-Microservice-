package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/gorilla/mux"
)

type URLpair struct {
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", indexFunc)
	router.HandleFunc("/new/", urlShortenerFunc)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}

func indexFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "hello")
}

func urlShortenerFunc(w http.ResponseWriter, r *http.Request) {
	userURL := r.URL.Path[5:]
	_, err := url.ParseRequestURI(userURL)
	if err != nil {
		log.Fatal(err)
	}
	response := URLpair{strings.Replace(userURL, "/", "//", 1), "localhost:8080/short" + "5"}

	byteData, err := json.Marshal(response)
	jsonResp := string(byteData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, jsonResp)
}
