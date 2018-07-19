package main

type AlexaRequest struct {
	Version string `json:"version"`
	Request struct {
		Type   string `json:"type"`
		Time   string `json:"timestamp"`
		Intent struct {
			Name               string `json:"name"`
			ConfirmationStatus string `json:"confirmationstatus"`
			Slots              map[string]struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

type AlexaTextResponse struct {
	Version  string `json:"version"`
	Response struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"outputSpeech"`
	} `json:"response"`
}

func NewAlexaTextResponse(text string) AlexaTextResponse {
	var resp AlexaTextResponse
	resp.Version = "1.0"
	resp.Response.OutputSpeech.Type = "PlainText"
	resp.Response.OutputSpeech.Text = text
	return resp
}
