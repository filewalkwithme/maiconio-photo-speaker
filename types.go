package main

type Photo struct {
	ID      int
	File    string
	Content string
	Labels  []Label
}

type Label struct {
	Name  string
	Score string
}

//ElephantsqlCredentials represents the credentials for Elephantsql
type ElephantsqlCredentials struct {
	URI string `json:"uri"`
}

//Elephantsql represents the Elephantsql connection
type Elephantsql struct {
	Credendials ElephantsqlCredentials `json:"credentials"`
}

//VisualRecognitionCredentials represents the credentials for Visual Recognition
type VisualRecognitionCredentials struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//VisualRecognition represents the VisualRecognition service
type VisualRecognition struct {
	Credendials VisualRecognitionCredentials `json:"credentials"`
}

//TextToSpeechCredentials represents the credentials for Text To Speech
type TextToSpeechCredentials struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//TextToSpeech represents the TextToSpeech service
type TextToSpeech struct {
	Credendials TextToSpeechCredentials `json:"credentials"`
}

//VCAPServices represents the VCAP environment var
type VCAPServices struct {
	Elephantsql       []Elephantsql       `json:"elephantsql"`
	VisualRecognition []VisualRecognition `json:"visual_recognition"`
	TextToSpeech      []TextToSpeech      `json:"text_to_speech"`
}

//WatsonResponse represents the response of the Watson Visual Recognition service
type WatsonResponse struct {
	Images []WatsonImage `json:"images"`
}

//WatsonImage represents the data about the image analyzed by Watson Visual Recognition service
type WatsonImage struct {
	Labels []WatsonLabel `json:"labels"`
}

//WatsonLabel represents the labels of the image analyzed by Watson Visual Recognition service
type WatsonLabel struct {
	Name  string `json:"label_name"`
	Score string `json:"label_score"`
}
