package telegram

import (
	"authService/config"
	client "authService/internal/grpc_clients"
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	bot        *tgbotapi.BotAPI
	userClient *client.UserClient
	bonusdb    *mongo.Database
)

func SetUserClient(client *client.UserClient) {
	userClient = client
}

func RunBot(cfg *config.Config) {
	var err error
	err = databaseConnection(cfg.Database.Port, cfg.Database.Host, cfg.Database.DbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bot, err = tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
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
			case update.Message.Text == "🎁 Bonus For Examify":
				handleBonusForExamify(update.Message)
			case update.Message.Text == "ℹ️ Details (Bonus)":
				handleDetailsOfBonus(update.Message)
			}
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}

func databaseConnection(databasePort, databaseHost, databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s/%s", databaseHost, databasePort, databaseName)
	clientOptions := options.Client().ApplyURI(uri)

	mongodbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = mongodbClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	bonusdb = mongodbClient.Database(databaseName)

	log.Printf("Connected to MongoDB database '%s' successfully", databaseName)
	return nil
}
