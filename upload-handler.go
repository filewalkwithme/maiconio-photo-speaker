package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

//soundHandler receives the "/upload" request
//this method sends the image to Watson Image Recognition Service and then gets
//the related labels (with scores)
func uploadHandler(w http.ResponseWriter, req *http.Request) {
	paramName := "img_File"

	//parse the from
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "req.ParseForm(): %v\n", err)
		return
	}

	//load the photo
	_, fileHeader, err := req.FormFile(paramName)
	if err != nil {
		fmt.Fprintf(w, "req.FormFile(\"img_File\"): %v\n", err)
		return
	}

	//open the photo file (to send the photo contents to Watson Image Recognition Service)
	file, err := fileHeader.Open()
	if err != nil {
		fmt.Fprintf(w, "fileHeader.Open(): %v\n", err)
		return
	}

	//creates a new form field (to be sent to Watson Image Recognition Service) containing
	//the received photo
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileHeader.Filename)
	if err != nil {
		fmt.Fprintf(w, "writer.CreateFormFile(paramName, fileHeader.Filename): %v\n", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Fprintf(w, "io.Copy(part, file): %v\n", err)
		return
	}

	//open a new photo file (this time to store on postgres)
	fileSave, err := fileHeader.Open()
	if err != nil {
		fmt.Fprintf(w, "[save]fileHeader.Open(): %v\n", err)
		return
	}

	fBuf, err := ioutil.ReadAll(fileSave)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll(file): %v\n", err)
		return
	}

	//convert to base64
	b64 := base64.StdEncoding.EncodeToString(fBuf)

	err = writer.Close()
	if err != nil {
		fmt.Fprintf(w, "writer.Close(): %v\n", err)
		return
	}

	//perform the request on Watson Image Recognition API
	request, err := http.NewRequest("POST", visualRecognitionURI, body)
	if err != nil {
		fmt.Fprintf(w, "http.NewRequest(\"POST\", visualRecognitionURI, body): %v\n", err)
		return
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(w, "client.Do(request): %v\n", err)
	} else {
		//get the response
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "body.ReadFrom(resp.Body): %v\n", err)
			return
		}

		err = resp.Body.Close()
		if err != nil {
			fmt.Fprintf(w, "resp.Body.Close(): %v\n", err)
			return
		}

		//store the response body on a string
		watsonResponseJSON := fmt.Sprintf("%v\n", body)

		//then deserialize it
		var watsonResponse WatsonResponse
		fmt.Printf("watsonResponseJSON: %v", watsonResponseJSON)
		err = json.Unmarshal([]byte(watsonResponseJSON), &watsonResponse)
		if err != nil {
			fmt.Fprintf(w, "json.Unmarshal(): %v\n", err)
			return
		}

		//inser the photo on postgres
		_, err = insertPhoto(fileHeader.Filename, b64, watsonResponse)
		if err != nil {
			fmt.Fprintf(w, "insertPhoto(fileHeader.Filename): %v\n", err)
			return
		}

		//redirect to index page
		http.Redirect(w, req, "/", 302)
	}
}
