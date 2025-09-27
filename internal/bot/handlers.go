package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "start":
		b.reply(update.Message.Chat.ID, "Добро пожаловать! Это OneCoffee Bot.")
	case "profile":
		b.reply(update.Message.Chat.ID, "Ваш профиль (пока пуст).")
	case "buy":
		b.reply(update.Message.Chat.ID, "Выберите подписку для покупки.")
	default:
		b.reply(update.Message.Chat.ID, "Неизвестная команда.")
	}
}

func (b *telegramBot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.tg.Send(msg)
}
