package handlers

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"tg_pager/internal/models"
	"tg_pager/internal/repo"
	"tg_pager/internal/services/ai"
	"tg_pager/internal/services/random"
	"tg_pager/internal/services/telegram"
	"tg_pager/internal/services/web"
)

const NameBot string = "Король"
const CommandSave = "с"
const CommandShow = "п"
const CommandRandom = "р"
const CommandSpeak = "г"
const LineBreak = "\r\n"

type Handler struct {
	botService    *telegram.Service
	webService    *web.WebServer
	aiService     *ai.AiService
	randomService *random.Service
	repo          *repo.Repository
}

func NewHandler(botService *telegram.Service, web *web.WebServer, repo *repo.Repository, rnd *random.Service, ai *ai.AiService) *Handler {
	return &Handler{
		botService:    botService,
		webService:    web,
		repo:          repo,
		randomService: rnd,
		aiService:     ai,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	msgChan := make(chan models.Message)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	go h.botService.RunBot(ctx, msgChan)
	err := h.webService.Start()
	if err != nil {
		cancel()
		return err
	}

	for {
		select {
		case msg := <-msgChan:
			h.checkMessage(msg)
		case <-sigChan:
			cancel()
			err := h.webService.Stop()
			if err != nil {
				return err
			}
			return nil
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

	if len(array) > 0 {
		if !strings.EqualFold(array[0], NameBot) {
			return
		}
	}
	if len(array) > 2 {
		// сохранить
		if strings.EqualFold(array[0], NameBot) && len(array[2]) > 0 {
			if strings.EqualFold(array[1], CommandSave) && len(msg.Msg) > 0 {
				h.saveMessage(array[2], msg.Msg)
				return
			}
		}
		// Показать
		if strings.EqualFold(array[1], CommandShow) {
			h.showMessage(msg, array[2])
			return
		}
		// запрос к ИИ
		if strings.EqualFold(array[0], NameBot) && strings.EqualFold(array[1], CommandSpeak) && len(array[2]) > 0 {
			prompt, ok := extractRemaining(msg.Command)
			if ok {
				h.showSpeak(msg, prompt)
			}
			return
		}

	}
	// Рандом
	if len(array) == 2 {
		if strings.EqualFold(array[0], NameBot) && strings.EqualFold(array[1], CommandRandom) {
			h.showRandom(msg)
			return
		}

	}
	// Помошь
	if len(array) == 1 {
		if strings.EqualFold(array[0], NameBot) {
			h.sendMessageError(msg)
		}
	}
}

func (h *Handler) sendMessageError(message models.Message) {
	answer := "Чего тебе? Я вот как работую: Укажи мое имя: " + NameBot + ", укажи команду: " + CommandSave + " (сохранить) или " + CommandShow + " (показать), укажи тэг. " +
		"Если хочешь что-то сохранить, ответь этой командой на сохраняемое сообщение. А если хочешь игрануть в 'Орел/Решка' просто набери " + NameBot + ", и укажи команду: " + CommandRandom + " (рандом)"
	h.botService.SendMessage(message.ID, answer)

}

func (h *Handler) getMessages(tag string) []string {
	m, err := h.repo.GetMessages(tag)
	if err != nil {
		return nil
	}
	messages := make([]string, 0, len(m))
	for i, val := range m {
		ms := strconv.Itoa(i) + ". " + val + LineBreak
		messages = append(messages, ms)
	}
	return messages
}

func (h *Handler) showMessage(in models.Message, tag string) {
	var sentence string
	out := h.getMessages(tag)

	if len(out) > 0 {
		sentence = strings.Join(out, "")
	} else {
		sentence = "НИ-ЧЕ-ГО!"
	}
	sentence = "Вот, что я нашел:" + LineBreak + sentence

	h.botService.SendMessage(in.ID, sentence)
}

func (h *Handler) random() bool {
	return h.randomService.GetRandom()
}

func (h *Handler) showRandom(in models.Message) {
	res := "Орел"
	coin := h.random()
	if coin {
		res = "Решка"
	}
	h.botService.SendMessage(in.ID, "Выпало: "+res)
}

func (h *Handler) showSpeak(in models.Message, prompt string) {
	response, err := h.aiService.GetResponse(prompt)
	if err != nil {
		response = " Не могу говорить, я занят! "

	}
	h.botService.SendMessage(in.ID, response)
}

func extractRemaining(str string) (string, bool) {
	word1Length := len(NameBot)
	word2Length := len(CommandSpeak)
	if len(str) < word1Length+1+word2Length {
		return "", false // Строка слишком короткая
	}
	remainingString := str[word1Length+1+word2Length:]
	return remainingString, len(remainingString) > 0
}
