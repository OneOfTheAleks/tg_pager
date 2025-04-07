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
