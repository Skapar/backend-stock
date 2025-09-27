package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	if update.Message.Contact != nil {
		phone := update.Message.Contact.PhoneNumber
		name := update.Message.Contact.FirstName + " " + update.Message.Contact.LastName
		tgUserID := update.Message.From.ID
		username := update.Message.From.UserName

		err := b.handleRegisterUser(tgUserID, username, name, phone)
		if err != nil {
			b.log.Error(err)
			b.reply(chatID, "Произошла ошибка при регистрации. Пожалуйста, попробуйте еще раз.")
			return
		}

		b.reply(chatID, "Спасибо! Мы получили ваш номер: "+phone+" (Имя: "+name+")")
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			b.handleStart(update)
		default:
			b.reply(chatID, "Неизвестная команда")
		}
	}
}

func (b *telegramBot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	send, err := b.tg.Send(msg)
	if err != nil {
		return
	}
	b.log.Infof("Sent message to %d: %s", chatID, send.Text)
}
