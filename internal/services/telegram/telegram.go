package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_pager/internal/models"
)

type Service struct {
	bot *tgbotapi.BotAPI
}

func New(token string) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &Service{bot: bot}, nil
}

func (s *Service) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (s *Service) RunBot(ctx context.Context, msgChan chan<- models.Message) {

	s.bot.Debug = true
	log.Printf("Authorized on account %s", s.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Bot stopping...")
			return
		case update := <-updates:
			if update.Message != nil {
				//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				//	msg.ReplyToMessageID = update.Message.MessageID

				//	_, err := s.bot.Send(msg)
				//	if err != nil {
				//		return
				//	}
				message := models.Message{
					Msg: update.Message.Text,
					ID:  update.Message.Chat.ID,
				}
				msgChan <- message
			}
		}
	}
}
