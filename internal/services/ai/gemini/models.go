package gemini

// Структуры для запроса к API Gemini
type GeminiRequest struct {
	Contents []Content `json:"contents"`
	// Можно добавить другие параметры, например, GenerationConfig
	// GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

/* // Опционально: Конфигурация генерации
type GenerationConfig struct {
	Temperature     float32  `json:"temperature,omitempty"`
	TopK            int      `json:"topK,omitempty"`
	TopP            float32  `json:"topP,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}
*/

// Структуры для ответа от API Gemini
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
	// Если API вернет ошибку в JSON теле при статусе 200 (маловероятно, но возможно)
	// или если статус не 200 и тело содержит ошибку
	Error *GeminiError `json:"error,omitempty"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// Структура для ошибок, возвращаемых API
type GeminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
