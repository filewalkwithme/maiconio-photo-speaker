package main

import (
	"fmt"
	"net/http"

	"text/template"
)

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

//indexHandler renders the main page
func indexHandler(w http.ResponseWriter, req *http.Request) {
	type IndexPage struct {
		AllPhotos []Photo
	}

	//gets all photos from postgres
	photos, err := allPhotos()
	if err != nil {
		fmt.Fprintf(w, "allPhotos(): %v\n", err)
		return
	}

	//then renders the main page
	var page = IndexPage{AllPhotos: photos}
	index.Execute(w, page)
}
