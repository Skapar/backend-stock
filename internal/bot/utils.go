package bot

func (b *telegramBot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	send, err := b.tg.Send(msg)
	if err != nil {
		return
	}
	b.log.Infof("Sent message to %d: %s", chatID, send.Text)
}
