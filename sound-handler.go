package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//soundHandler receives the "/sound" request
//it is a proxy to Watson TextToSpeech API
func soundHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	text := params.Get("text")
	text = url.QueryEscape(text)

	//calls Watson TextToSpeech API
	fmt.Printf("soundHandler: %v \n", textToSpeechURI+"?text="+text)
	request, _ := http.NewRequest("GET", textToSpeechURI+"?text="+text, nil)
	copyHeader(req.Header, &request.Header)

	client := &http.Client{}
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//Forward the response
	w.Write(body)
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
