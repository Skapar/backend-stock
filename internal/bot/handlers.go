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
			return
		}
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case STARTBUTTONNAME:
			b.handleStart(update)
		default:
			b.reply(chatID, UNEXISTINGBUTTONPRESSED)
		}
	}
}

func (b *telegramBot) handleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	contactButton := tgbotapi.NewKeyboardButtonContact(SENDCONTACTTEXT)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, STARTGREETINGTEXT)
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
		b.reply(tgUserID, ERRORREGISTRATIONTEXT)
		b.log.Errorf("failed to register user: %v", err)
		return err
	}

	b.reply(tgUserID, fmt.Sprintf(THANKYOUREGISTERTEXT, name))
	return nil
}
