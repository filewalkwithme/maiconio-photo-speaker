package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var vcap VCAPServices
var visualRecognitionURI string
var textToSpeechURI string

//initAPP check if all environment variables are configured and try to connect
//with postgres
func initAPP() {
	vcapJSON := os.Getenv("VCAP_SERVICES")

	//Check the presence of VCAP_SERVICES
	if len(vcapJSON) == 0 {
		fmt.Printf("VCAP_SERVICES not defined\n")
		os.Exit(1)
	}

	//load the environment variable
	err := json.Unmarshal([]byte(vcapJSON), &vcap)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	//check the presense of VisualRecognition service
	if len(vcap.VisualRecognition) == 0 {
		fmt.Printf("VisualRecognition not defined on VCAP_SERVICES\n")
		os.Exit(1)
	}

	visualRecognitionURI = vcap.VisualRecognition[0].Credendials.URL
	visualRecognitionURI = strings.Replace(visualRecognitionURI, "https://", "https://"+vcap.VisualRecognition[0].Credendials.Username+":"+vcap.VisualRecognition[0].Credendials.Password+"@", -1)
	visualRecognitionURI = visualRecognitionURI + "/v1/tag/recognize"
	fmt.Printf("visualRecognitionURI: %v\n", visualRecognitionURI)

	//check the presense of TextToSpeech service
	if len(vcap.TextToSpeech) == 0 {
		fmt.Printf("TextToSpeech not defined on VCAP_SERVICES\n")
		os.Exit(1)
	}

	textToSpeechURI = vcap.TextToSpeech[0].Credendials.URL
	textToSpeechURI = strings.Replace(textToSpeechURI, "https://", "https://"+vcap.TextToSpeech[0].Credendials.Username+":"+vcap.TextToSpeech[0].Credendials.Password+"@", -1)
	textToSpeechURI = textToSpeechURI + "/v1/synthesize"
	fmt.Printf("textToSpeechURI: %v\n", textToSpeechURI)

	//check the presense of Elephantsql service
	if len(vcap.Elephantsql) == 0 {
		fmt.Printf("Elephantsql not defined on VCAP_SERVICES\n")
		os.Exit(1)
	}

	//Try to open a connection to postgres
	fmt.Printf("Connecting on %v\n", vcap.Elephantsql[0].Credendials.URI)
	tmpDB, err := sql.Open("postgres", vcap.Elephantsql[0].Credendials.URI)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	db = tmpDB

	//Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DB Connected! \n")
}
