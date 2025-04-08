package models

// сообщение
type Message struct {
	Msg       string
	ChatID    int64
	Command   string
	MessageID int
}
