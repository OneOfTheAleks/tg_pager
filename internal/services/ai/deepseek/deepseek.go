package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DSService struct {
	apiKey string
}

func New(apiKey string) *DSService {
	return &DSService{apiKey: apiKey}
}

func (s *DSService) GetResponse(prompt string) (string, error) {
	m := Message{
		Role:    "user",
		Content: prompt,
	}
	mArray := []*Message{&m}
	reqBody := ChatCompletionsRequest{
		Model:    "deepseek-chat",
		Messages: mArray,
	}

	jsonData, err := json.Marshal(&reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Обработка ответа (пример)

	choicesRaw, ok := result["choices"]
	if !ok || choicesRaw == nil {
		return "", fmt.Errorf("response does not contain 'choices' or it is nil")
	}

	choices, ok := result["choices"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected type for 'choices': expected []interface{}, got %T", choicesRaw)
	}

	if len(choices) > 0 {
		message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
		return message["content"].(string), nil
	}

	return "", fmt.Errorf("no response from DeepSeek")
}
