package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case STARTBUTTON:
		b.sendStartMessage(update.Message.Chat.ID)
	case GETPAYMENTDETAILSBUTTON:
		b.sendPaymentDetails(update.Message.Chat.ID)
	case ASKRECEIPTBUTTON:
		b.askForReceipt(update.Message.Chat.ID)
	default:
		err := b.reply(update.Message.Chat.ID, UNEXISTINGBUTTONPRESSED)
		if err != nil {
			b.log.Errorf("failed to send unexisting button message: %v", err)
		}
	}
}

func (b *telegramBot) handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {
	case GETPAYMENTDETAILSBUTTON:
		b.sendPaymentDetails(chatID)
	case ASKRECEIPTBUTTON:
		b.askForReceipt(chatID)
	default:
		b.reply(chatID, UNEXISTINGBUTTONPRESSED)
	}

	resp, err := b.tg.Request(tgbotapi.NewCallback(callback.ID, ""))
	if err != nil {
		b.log.Errorf("failed to answer callback query: %v", err)
	}
	if !resp.Ok {
		b.log.Errorf("callback query response error: %v", resp.Description)
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
	getPaymentBtn := tgbotapi.NewInlineKeyboardButtonData(GETPAYMENTDETAILSBUTTONTEXT, GETPAYMENTDETAILSBUTTON)
	sendReceiptBtn := tgbotapi.NewInlineKeyboardButtonData(ASKRECEIPTBUTTONTEXT, ASKRECEIPTBUTTON)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(getPaymentBtn),
		tgbotapi.NewInlineKeyboardRow(sendReceiptBtn),
	)

	msg := tgbotapi.NewMessage(chatID, POSTREGISTRATIONTEXT)
	msg.ReplyMarkup = keyboard

	b.sendMessage(msg)
}

func (b *telegramBot) sendPaymentDetails(chatID int64) {
	err := b.reply(chatID, PAYMENTDETAILSTEXT)
	if err != nil {
		b.log.Errorf("failed to send payment details: %v", err)
	}
}

func (b *telegramBot) askForReceipt(chatID int64) {
	err := b.reply(chatID, ASKRECEIPTTEXT)
	if err != nil {
		b.log.Errorf("failed to ask for receipt: %v", err)
	}
}
