package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	//DOMAIN const for domain name
	DOMAIN = "http://localhost"
	// PORT const for port value
	PORT = ":8080"
)

// URLpair struct for json data to be generated
type URLpair struct {
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

func main() {
	router := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	router.Handle("/static/", http.StripPrefix("/static/", files))

	router.HandleFunc("/", indexFunc)
	router.HandleFunc("/new/", urlShortenerFunc)
	router.HandleFunc("/submit/", submitFunc)
	router.HandleFunc("/redirect/", redirectToOriginal)

	server := http.Server{
		Addr:    PORT,
		Handler: router,
	}
	server.ListenAndServe()
}

// indexFunc display the home page if request is sent to the root
// of the app otherwise call redirectToOriginal function
func indexFunc(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("template/index.html"))
	if r.URL.Path != "/" {
		redirectToOriginal(w, r)
		return
	}
	templ.Execute(w, nil)
}

// submitFunc handle request coming from form submition from home page
// make get the long_url value, create the short url and give back the json data
func submitFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rawURL := r.Form["long_url"][0]
	newURLPair, err := createShortURL(rawURL)
	if err != nil {
		displayErrorJSON(w, rawURL)
		return
	}
	if err := newURLPair.getUrls(); err != nil {
		fmt.Fprint(w, generateResponse(newURLPair))
		return
	}
	fmt.Fprint(w, generateResponse(newURLPair))
}

//This function try to find short url in the databse
// and then redirect to the original long url
func redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	rURL := URLpair{}
	rURL.Short = DOMAIN + PORT + r.URL.RequestURI()
	statement := "SELECT ORIGINAL FROM URLMAP WHERE short_url=$1"
	err := db.QueryRow(statement, rURL.Short).Scan(&rURL.Original)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	http.Redirect(w, r, rURL.Original, 303)
}

// urlShortenestFunc take long url from url field of the browser
// create the short url corresponding and send back the json data
func urlShortenerFunc(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.RequestURI()[5:]
	newURLPair, err := createShortURL(rawURL)
	if err != nil {
		displayErrorJSON(w, rawURL)
		return
	}

	if err := newURLPair.getUrls(); err != nil {
		fmt.Fprint(w, generateResponse(newURLPair))
		return
	}
	fmt.Fprint(w, generateResponse(newURLPair))

}

// this fonction take a url string parse it to check if it is valid
// then create a short url of it
func createShortURL(URL string) (URLpair, error) {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return URLpair{}, err
	}
	newURL := URLpair{URL, DOMAIN + PORT + "/redirect/" + getURLCount() + "/"}

	return newURL, nil
}

// this function take a URLpair data and make a json string from it
func generateResponse(newU URLpair) string {
	byteData, err := json.Marshal(newU)
	jsonResp := string(byteData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonResp
}

// this function return the number of records in the database incremented by 1
func getURLCount() string {
	var idNumber int
	_ = db.QueryRow("SELECT COUNT(*) FROM URLMAP").Scan(&idNumber)
	return strconv.Itoa(idNumber + 1)
}

func displayErrorJSON(w http.ResponseWriter, originalURL string) {
	errorResponse := URLpair{originalURL, "error invalid url format!"}
	fmt.Fprint(w, generateResponse(errorResponse))
}
