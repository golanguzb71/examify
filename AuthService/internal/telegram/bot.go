package telegram

import (
	"authService/config"
	client "authService/internal/grpc_clients"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot        *tgbotapi.BotAPI
	userClient *client.UserClient
)

func SetUserClient(client *client.UserClient) {
	userClient = client
}

func RunBot(cfg *config.Config) {
	var err error
	bot, err = tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch {
			case update.Message.IsCommand():
				handleCommand(update.Message)
			case update.Message.Contact != nil:
				handleContact(update.Message)
			}
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}
