package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HugService struct {
	apiKey string
}

func NewHug(apiKey string) *HugService {
	return &HugService{apiKey: apiKey}
}

// https://api-inference.huggingface.co/models/DeepSeek-V3-0324

func (h *HugService) GetResponse(prompt string) (string, error) {
	// Проверяем, что prompt не пустой
	if prompt == "" {
		fmt.Println("Prompt не может быть пустым")
		return "", errors.New("empty prompt")
	}

	// Формируем запрос
	requestBody := &HuggingFaceRequest{
		Inputs: prompt,
		Parameters: struct {
			MaxNewTokens int `json:"max_new_tokens,omitempty"`
		}{
			MaxNewTokens: 100, // Опционально: ограничиваем длину ответа
		},
	}

	// Маршалим JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Ошибка при маршалинге JSON: %v\n", err)
		return "", err
	}

	// Проверяем, что jsonData не пустой
	if len(jsonData) == 0 {
		fmt.Println("jsonData пуст")
		return "", errors.New("empty JSON data")
	}

	// Логируем JSON для отладки
	fmt.Printf("Отправляемый JSON: %s\n", string(jsonData))

	// Создаем HTTP-запрос
	buf := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest(
		"POST",
		"https://api-inference.huggingface.co/models/deepseek-ai/DeepSeek-V3",
		buf,
	)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return "", err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return "", err
	}

	// Разбираем JSON
	var apiResponse HuggingFaceResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return "", err
	}

	// Выводим ответ
	answer := "Пустой ответ от API"
	if len(apiResponse) > 0 {
		answer = apiResponse[0].GeneratedText
	}

	return answer, nil
}
