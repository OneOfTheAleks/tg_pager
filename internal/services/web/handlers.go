package web

import (
	"github.com/gofiber/fiber/v2"
	"tg_pager/internal/repo"
)

type WebHandlers struct {
	repo *repo.Repository
}

func newWebHandlers(repo *repo.Repository) *WebHandlers {
	return &WebHandlers{
		repo: repo,
	}
}

func (w *WebHandlers) Tags(c *fiber.Ctx) error {
	tags, err := w.repo.GetTags()
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Теги не найдены")
	}

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

	res, err := h.repo.GetMessages(tag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Контент по тегу не найден")
	}

	if len(res) < 1 {
		return c.Status(fiber.StatusNotFound).SendString("Контент по тегу не найден")
	}

	// Подготавливаем структуру с индексами
	type ContentRow struct {
		Index   int
		Content string
	}

	var rows []ContentRow
	for i, text := range res {
		rows = append(rows, ContentRow{
			Index:   i + 1, // начинаем с 1
			Content: text,
		})
	}

	return c.Render("content", fiber.Map{
		"Title":   "Контент для тега: " + tag,
		"Tag":     tag,
		"Content": rows,
	})

}
