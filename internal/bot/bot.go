package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/onec-tech/bot/internal/models/entities"
	"github.com/onec-tech/bot/internal/service"
	"github.com/onec-tech/bot/pkg/logger"
)

type telegramBot struct {
	tg      *tgbotapi.BotAPI
	log     logger.Logger
	service service.Service
}

func NewBot(token string, log logger.Logger, service service.Service) (Bot, error) {
	tg, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	tg.Debug = false

	return &telegramBot{
		tg:      tg,
		log:     log,
		service: service,
	}, nil
}

func (b *telegramBot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tg.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(update)
		case <-ctx.Done():
			return ctx.Err()
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
