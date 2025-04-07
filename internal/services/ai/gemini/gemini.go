package gemini

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Gemini struct {
	apiKey    string
	modelName string
}

func New(apikey, modelName string) (*Gemini, error) {
	if modelName == "" {
		modelName = "gemini-2.0-flash" // Модель по умолчанию
	}
	if apikey == "" {
		return nil, errors.New("Не указан api key")
	}
	return &Gemini{
		apiKey:    apikey,
		modelName: modelName,
	}, nil
}

func (r *Gemini) GetResponse(prompt string) (string, error) {

	prompt = "Представь что ты Саймон Петриков Снежный король из мультфильма Время Приключений.  Как бы ты ответил на это: " + prompt

	// Формируем URL эндпоинта
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", r.modelName, r.apiKey)

	// Формируем тело запроса
	requestBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
		// Можно раскомментировать и настроить GenerationConfig здесь, если нужно
		/*
			GenerationConfig: &GenerationConfig{
				Temperature:     0.7,
				MaxOutputTokens: 1000,
			},
		*/
	}

	// Кодируем тело запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодирования JSON запроса: %w", err)
	}

	// Создаем HTTP POST запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка создания HTTP запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения HTTP запроса: %w", err)
	}
	defer resp.Body.Close() // Важно закрыть тело ответа

	// Читаем тело ответа
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		// Пытаемся декодировать ошибку из тела ответа
		var apiError GeminiError
		if json.Unmarshal(responseBodyBytes, &apiError) == nil && apiError.Message != "" {
			// Если удалось распарсить структуру ошибки Gemini
			return "", fmt.Errorf("ошибка API Gemini (статус %d): %s (%s)", resp.StatusCode, apiError.Message, apiError.Status)
		}
		// Если не удалось распарсить как GeminiError, возвращаем как есть
		return "", fmt.Errorf("ошибка API (статус %d): %s", resp.StatusCode, string(responseBodyBytes))
	}

	// Декодируем успешный ответ JSON
	var geminiResponse GeminiResponse
	err = json.Unmarshal(responseBodyBytes, &geminiResponse)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON ответа: %w. Тело ответа: %s", err, string(responseBodyBytes))
	}

	// Проверяем, есть ли кандидаты в ответе
	if len(geminiResponse.Candidates) == 0 {
		// Проверим, не вернул ли API ошибку в поле error при статусе 200
		if geminiResponse.Error != nil {
			return "", fmt.Errorf("API Gemini вернуло ошибку: %s (%s)", geminiResponse.Error.Message, geminiResponse.Error.Status)
		}
		return "", fmt.Errorf("API не вернуло кандидатов в ответе. Тело ответа: %s", string(responseBodyBytes))
	}

	// Проверяем, есть ли части контента в первом кандидате
	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("API не вернуло текстовых частей в ответе кандидата. Тело ответа: %s", string(responseBodyBytes))
	}

	// Возвращаем текст из первой части первого кандидата
	return geminiResponse.Candidates[0].Content.Parts[0].Text, nil
}
