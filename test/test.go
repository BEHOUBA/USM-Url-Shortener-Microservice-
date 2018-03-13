package main

import (
	"fmt"
	"net/url"
)

func checkURL(urlString string) bool {
	if urlString[:8] == "https://" || urlString[:7] == "http://" {
		_, err := url.ParseRequestURI(urlString)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func main() {
	fmt.Println(checkURL("www.google.com"))
}
