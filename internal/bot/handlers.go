package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/onec-tech/bot/internal/models/entities"
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

func (b *telegramBot) handleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// Текст приветствия
	text := "Добро пожаловать в OneCoffee! Пожалуйста, отправьте ваше имя и номер телефона."

	// Кнопка для отправки контакта
	contactButton := tgbotapi.NewKeyboardButtonContact("Отправить свой номер")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err := b.tg.Send(msg)
	if err != nil {
		b.log.Errorf("failed to send start message: %v", err)
	}
}

func (b *telegramBot) handleRegisterUser(tgUserID int64, username, name, phone string) error {
	user := &entities.User{
		TGID:     tgUserID,
		Username: username,
		Name:     name,
		Phone:    phone,
	}

	err := b.service.CreateOrUpdateUser(context.Background(), user)
	if err != nil {
		b.reply(tgUserID, "Произошла ошибка при регистрации, попробуйте позже.")
		b.log.Errorf("failed to register user: %v", err)
		return err
	}

	b.reply(tgUserID, fmt.Sprintf("Спасибо, %s! Вы успешно зарегистрированы в OneCoffee.", name))
	return nil
}
