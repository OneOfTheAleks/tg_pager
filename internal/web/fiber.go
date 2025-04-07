package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	addr string
	port string
	app  *fiber.App
}

func New(addr string, port string) (*WebServer, error) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	return &WebServer{
		addr: addr,
		port: port,
		app:  app,
	}, nil
}

func (w *WebServer) Start() error {
	// Создаём канал для ошибок
	errChan := make(chan error, 1)

	// Запускаем сервер в горутине
	go func() {
		if err := w.app.Listen(w.addr + ":" + w.port); err != nil {
			// Логируем ошибку и отправляем её в канал
			fmt.Printf("Ошибка при запуске сервера: %v\n", err)
			errChan <- fmt.Errorf("ошибка при запуске сервера: %v", err)
			return
		}
		close(errChan) // Закрываем канал, если ошибок нет
	}()

	// Не ждём ошибку из канала, а сразу возвращаем управление
	return nil
}

func (w *WebServer) Stop() error {
	return w.app.Shutdown()
}
