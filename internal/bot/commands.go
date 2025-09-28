package bot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case STARTBUTTON:
		b.sendStartMessage(update.Message.Chat.ID)
	case "mysubscription":
		b.sendMySubscription(update.Message.Chat.ID, update.Message.From.ID)
	default:
		b.reply(update.Message.Chat.ID, UNEXISTINGBUTTONPRESSED)
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

	inlineBtn := tgbotapi.NewInlineKeyboardButtonData(MYSUBSCRIPTIONTEXT, "mysubscription")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(inlineBtn),
	)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, STARTGREETINGTEXT)
	msg.ReplyMarkup = replyKeyboard
	b.sendMessage(msg)

	// Отправляем вторым сообщением кнопку подписки
	btnMsg := tgbotapi.NewMessage(chatID, "Ваши действия:")
	btnMsg.ReplyMarkup = inlineKeyboard
	b.sendMessage(btnMsg)
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
	b.reply(chatID, PAYMENTDETAILSTEXT)
}

func (b *telegramBot) askForReceipt(chatID int64) {
	b.reply(chatID, ASKRECEIPTTEXT)
}

func (b *telegramBot) handleReceiptDocument(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	tgID := update.Message.From.ID
	doc := update.Message.Document

	user, err := b.service.GetUserByTGID(context.Background(), tgID)
	if err != nil || user == nil {
		b.reply(chatID, "Не удалось найти пользователя в системе.")
		b.log.Errorf("user not found for tgID %d: %v", tgID, err)
		return
	}

	// Создаём директорию, если нет
	if err := os.MkdirAll("./receipts", os.ModePerm); err != nil {
		b.reply(chatID, "Не удалось создать директорию для файлов.")
		b.log.Errorf("failed to create receipts dir: %v", err)
		return
	}

	// Скачиваем файл
	fileID := doc.FileID
	file, err := b.tg.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		b.reply(chatID, "Не удалось получить файл. Попробуйте снова.")
		b.log.Errorf("failed to get file: %v", err)
		return
	}

	url := file.Link(b.tg.Token)
	filePath := fmt.Sprintf("./receipts/%d_%s", tgID, doc.FileName)

	out, err := os.Create(filePath)
	if err != nil {
		b.reply(chatID, "Не удалось создать файл для сохранения.")
		b.log.Errorf("failed to create file: %v", err)
		return
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		b.reply(chatID, "Не удалось скачать файл.")
		b.log.Errorf("failed to download file: %v", err)
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		b.reply(chatID, "Не удалось сохранить файл.")
		b.log.Errorf("failed to copy file: %v", err)
		return
	}

	// Сохраняем в базу, используя внутренний user.ID
	if err := b.service.CreateReceipt(context.Background(), user.ID, filePath); err != nil {
		b.reply(chatID, "Не удалось сохранить информацию о чеке.")
		b.log.Errorf("failed to save receipt: %v", err)
		return
	}

	b.reply(chatID, "Файл успешно сохранён. Спасибо!")
}

func (b *telegramBot) sendMySubscription(chatID, tgID int64) {
	ctx := context.Background()

	user, err := b.service.GetUserByTGID(ctx, tgID)
	if err != nil || user == nil {
		b.reply(chatID, "Не удалось найти пользователя.")
		return
	}

	sub, err := b.service.GetActiveSubscription(ctx, user.ID)
	if err != nil || sub == nil {
		b.reply(chatID, "У вас пока нет активной подписки.")
		return
	}

	msg := fmt.Sprintf(
		"Ваша подписка:\nНачало: %s\nОкончание: %s\nОсталось чашек: %d из %d",
		sub.StartDate.Format("02.01.2006"),
		sub.EndDate.Format("02.01.2006"),
		sub.RemainingCups,
		sub.TotalCups,
	)

	b.reply(chatID, msg)
}
