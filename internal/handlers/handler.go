package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tg_pager/internal/models"
	"tg_pager/internal/repo"
	"tg_pager/internal/services/telegram"
)

const NameBot string = "Король"
const CommandSave = "с"
const CommandShow = "п"

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
			h.checkMessage(msg)
		case <-sigChan:
			cancel()
			return
		}
	}

}

func (h *Handler) saveMessage(tag, msg string) {

	err := h.repo.SaveMessage(tag, msg)
	if err != nil {
		return
	}
}

func (h *Handler) checkMessage(msg models.Message) {
	array := strings.Fields(msg.Command)
	println(array[0])
	println(NameBot)
	println(strings.EqualFold(array[0], NameBot))
	if len(array) > 0 {
		if !strings.EqualFold(array[0], NameBot) {
			return
		}
	}
	if len(array) > 2 {
		if strings.EqualFold(array[0], NameBot) || strings.EqualFold(array[1], "с") || len(array[2]) > 0 || len(msg.Msg) > 0 {
			h.saveMessage(array[2], msg.Msg)
			return
		}
	}
	if strings.EqualFold(array[0], NameBot) {
		h.sendMessageError(msg)
	}

}

func (h *Handler) sendMessageError(message models.Message) {
	answer := "Ну нет, я вот как работую: Укажи мое имя: " + NameBot + ", укажи команду: " + CommandSave + " (сохранить) или " + CommandShow + " (показать), укажи тэг. Если хочешь что-то сохранить, ответь этой командойн на сохраняемое сообщение"
	h.botService.SendMessage(message.ID, answer)

}
