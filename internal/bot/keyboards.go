package bot

import (
	"fmt"
	"github.com/onec-tech/bot/internal/models/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func buildSubscriptionsKeyboard(subs []entities.Subscription) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, s := range subs {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			s.Name+" – "+formatPrice(s.Price),
			fmt.Sprintf("buy_%d", s.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func formatPrice(price int64) string {
	return fmt.Sprintf("%d₸", price)
}
