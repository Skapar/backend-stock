package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/onec-tech/bot/internal/models/entities"
)

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

func (b *telegramBot) registerUserFromContact(update tgbotapi.Update) {
	phone := update.Message.Contact.PhoneNumber
	name := update.Message.Contact.FirstName + " " + update.Message.Contact.LastName
	tgUserID := update.Message.From.ID
	username := update.Message.From.UserName

	if err := b.handleRegisterUser(tgUserID, username, name, phone); err != nil {
		b.log.Error(err)
		return
	}

	b.sendPostRegistrationKeyboard(tgUserID)
}
