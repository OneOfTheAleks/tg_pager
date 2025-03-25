package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tg_pager/internal/models"
	"tg_pager/internal/repo"
	"tg_pager/internal/services/telegram"
)

type Handler struct {
	botService *telegram.Service
	//	deepseekService *deepseek.Service
	//	randomService   *random.Service
	repo *repo.Repository
}

func NewHandler(botService *telegram.Service, repo *repo.Repository) *Handler {
	return &Handler{
		botService: botService,
		repo:       repo,
	}
}

func (h *Handler) Start(ctx context.Context) {
	msgChan := make(chan models.Message)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	go h.botService.RunBot(ctx, msgChan)

	for {
		select {
		case msg := <-msgChan:
			fmt.Println("Получено:", msg)
		case <-sigChan:
			cancel()
			return
		}
	}

}
