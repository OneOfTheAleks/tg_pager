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
	contentMap := map[string]string{
		"golang":          "Golang — это язык программирования, разработанный Google.",
		"fiber":           "Fiber — это веб-фреймворк для Go, вдохновленный Express.js.",
		"web-development": "Веб-разработка — это создание веб-приложений и сайтов.",
	}

	if content, exists := contentMap[tag]; exists {
		return c.Render("content", fiber.Map{
			"Title":   "Контент для тега: " + tag,
			"Tag":     tag,
			"Content": content,
		})
	}

	return c.Status(fiber.StatusNotFound).SendString("Тег не найден")
}
