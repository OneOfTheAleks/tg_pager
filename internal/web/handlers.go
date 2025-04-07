package web

import "github.com/gofiber/fiber/v2"

type WebHandlers struct {
}

func newWebHandlers() *WebHandlers {
	return &WebHandlers{}
}

func (w *WebHandlers) Tags(c *fiber.Ctx) error {
	tags := []string{"golang", "fiber", "web-development"}
	return c.Render("tags", fiber.Map{
		"Title": "Список тегов",
		"Tags":  tags,
	})
}

func (w *WebHandlers) Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Главная страница",
	})
}

func (h *WebHandlers) Content(c *fiber.Ctx) error {
	tag := c.Query("tag")

	// Структура данных: тег -> список записей
	contentMap := map[string][]string{
		"golang": {
			"Golang — это язык программирования, разработанный Google.",
			"Go используется для создания высокопроизводительных приложений.",
			"Go имеет встроенную поддержку конкурентности через горутины.",
		},
		"fiber": {
			"Fiber — это веб-фреймворк для Go, вдохновленный Express.js.",
			"Fiber обеспечивает высокую производительность и простоту использования.",
		},
		"web-development": {
			"Веб-разработка — это создание веб-приложений и сайтов.",
			"Современная веб-разработка включает фронтенд и бэкенд.",
		},
	}

	// Проверяем, существует ли тег
	if contentList, exists := contentMap[tag]; exists {
		return c.Render("content", fiber.Map{
			"Title":   "Контент для тега: " + tag,
			"Tag":     tag,
			"Content": contentList,
		})
	}

	// Если тег не найден, возвращаем ошибку
	return c.Status(fiber.StatusNotFound).SendString("Тег не найден")
}
