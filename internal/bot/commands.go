package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case STARTBUTTONNAME:
		b.sendStartMessage(update.Message.Chat.ID)
	default:
		b.reply(update.Message.Chat.ID, UNEXISTINGBUTTONPRESSED)
	}
}

func (b *telegramBot) sendStartMessage(chatID int64) {
	contactButton := tgbotapi.NewKeyboardButtonContact(SENDCONTACTTEXT)
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(contactButton))
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, STARTGREETINGTEXT)
	msg.ReplyMarkup = keyboard
	b.sendMessage(msg)
}

func (b *telegramBot) sendPostRegistrationKeyboard(chatID int64) {
	getPaymentBtn := tgbotapi.NewKeyboardButton(GETPAYMENTTEXT)
	sendReceiptBtn := tgbotapi.NewKeyboardButton(SENDRECEIPTTEXT)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPaymentBtn),
		tgbotapi.NewKeyboardButtonRow(sendReceiptBtn),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	b.sendMessage(msg)
}
