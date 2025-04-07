package deepseek

// Структуры для запроса и ответа API
type HuggingFaceRequest struct {
	Inputs     string `json:"inputs"`
	Parameters struct {
		MaxNewTokens int `json:"max_new_tokens,omitempty"`
	} `json:"parameters,omitempty"`
}

type HuggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}
