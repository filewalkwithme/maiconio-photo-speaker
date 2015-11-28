package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

//deleteHandler receives the "/delete" request
func deleteHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "strconv.Atoi(idStr): %v\n", err)
			return
		}

		//remove the photo from postgres
		n, err := removePhoto(id)
		if err != nil {
			fmt.Fprintf(w, "removePhoto(id): %v\n", err)
			return
		}

		//and remove from disk, too
		err = os.Remove("static/images/" + idStr + ".jpg")
		if err != nil {
			fmt.Fprintf(w, "os.Remove(%v): %v\n", idStr+".jpg", err)
			return
		}

		fmt.Printf("Rows removed: %v\n", n)
	}
	http.Redirect(w, req, "/", 302)
}
