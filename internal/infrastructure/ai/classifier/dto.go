package classifier

type ClassifyRequest struct {
	Text string `json:"text"`
}

type ClassifyResponse struct {
	Level      string         `json:"level"`
	Confidence float64        `json:"confidence"`
	Metadata   map[string]any `json:"metadata"`
}
