package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		sendWelcome(message)
	}
}

func sendWelcome(message *tgbotapi.Message) {
	name := message.From.FirstName

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📞 Send Contact"),
		),
	)

	welcomeMessage := fmt.Sprintf(`
🇺🇿
Salom %s 👋
CodeVan servicega xush kelibsiz
⬇️ Kontaktingizni yuboring va 10 daqiqalik kalitingizni oling!

🇺🇸
Hi %s 👋
Welcome to CodeVan service
⬇️ Send your contact and get 10 minutes key!
`, name, name)

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeMessage)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handleContact(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	contact := message.Contact
	phoneNumber := contact.PhoneNumber
	firstName := contact.FirstName
	lastName := contact.LastName

	if userClient == nil {
		log.Println("Error: userClient is nil")
		return
	}
	user, err := userClient.GetUserByChatIdOrPhone(nil, &phoneNumber, nil)

	if err != nil || user == nil {
		if lastName == "" {
			lastName = "********"
		}
		_, err = userClient.CreateUser(firstName, lastName, fmt.Sprintf("%v", chatID), phoneNumber)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return
		}
	}

	//requiredChannels := []string{"mock_examify"}
	//isMember, err := checkExistInRequiredChannel(chatID, requiredChannels)
	//if err != nil {
	//	log.Printf("Error checking channel membership: %v", err)
	//	return
	//}

	//if !isMember {
	//	msg := tgbotapi.NewMessage(chatID, "Please join our required channels before proceeding. @mock_examify")
	//	bot.Send(msg)
	//	return
	//}

	existingCode := GetStoredCode(fmt.Sprintf("%v", chatID))
	if existingCode != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Eski kodingiz hali ham kuchda ☝️ <code>%s</code>", *existingCode))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}
	code := generateCode()
	storeCode(code, chatID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Yangilash / Renew", fmt.Sprintf("renew_%s_%d", code, chatID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔒 Code: <code>%s</code>", code))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	if strings.HasPrefix(callback.Data, "renew_") {
		parts := strings.Split(callback.Data, "_")
		if len(parts) != 3 {
			return
		}
		chatID := callback.Message.Chat.ID

		existingCode := GetStoredCode(fmt.Sprintf("%v", chatID))
		if existingCode != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Eski kodingiz hali ham kuchda ☝️ <code>%s</code>", *existingCode))
			msg.ParseMode = "HTML"
			bot.Send(msg)
			return
		}

		newCode := generateCode()
		storeCode(newCode, chatID)

		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, fmt.Sprintf("🔒 Code: <code>%s</code>", newCode))
		edit.ParseMode = "HTML"
		edit.ReplyMarkup = callback.Message.ReplyMarkup
		bot.Send(edit)
	}
}
