package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	send, err := b.tg.Send(msg)
	if err != nil {
		return
	}
	b.log.Infof("Sent message to %d: %s", chatID, send.Text)
}

func (b *telegramBot) sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := b.tg.Send(msg); err != nil {
		b.log.Errorf("failed to send message: %v", err)
	}
}
