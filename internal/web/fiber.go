package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type WebServer struct {
	addr string
	port string
	app  *fiber.App
}

func New(addr string, port string) (*WebServer, error) {

	engine := html.New("./html", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	h := newWebHandlers()
	app.Get("/", h.Home)
	app.Get("/tags", h.Tags)
	app.Get("/content", h.Content)

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
