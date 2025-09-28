package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		b.handleCallback(update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	if update.Message.Document != nil {
		b.handleReceiptDocument(update)
		return
	}

	if update.Message.Contact != nil {
		b.registerUserFromContact(update)
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(update)
		return
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
