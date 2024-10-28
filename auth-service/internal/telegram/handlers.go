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
		referralCode := message.CommandArguments()
		if referralCode != "" {
			err := processReferral(message.Chat.ID, referralCode, message.From.FirstName)
			if err != nil {
				log.Printf("Error processing referral: %v", err)
			}
		}
		sendWelcome(message)
	}
}

func sendWelcome(message *tgbotapi.Message) {
	name := message.From.FirstName

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📞 Send Contact"),
			tgbotapi.NewKeyboardButton("🎁 Bonus For Examify"),
			tgbotapi.NewKeyboardButton("ℹ️ Details (Bonus)"),
		),
	)

	welcomeMessage := fmt.Sprintf(`
🇺🇿
Salom %s 👋
CodeVan servicega xush kelibsiz
examify.uz ro'yxatdan o'tish uchun ⬇️ kontaktingizni yuboring va 1 daqiqalik kalitingizni oling!

🇺🇸
Hi %s 👋
Welcome to CodeVan service
⬇️Send your contact and get 1 minutes key!
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
	if phoneNumber[0] != '+' {
		phoneNumber = "+" + phoneNumber
	}
	user, err := userClient.GetUserByChatIdOrPhone(nil, &phoneNumber, nil)

	if err != nil && user == nil {
		if lastName == "" {
			lastName = "********"
		}
		_, err = userClient.CreateUser(firstName, lastName, fmt.Sprintf("%v", chatID), phoneNumber)
		if err != nil {
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
	chatID := callback.Message.Chat.ID
	if strings.HasPrefix(callback.Data, "renew_") {
		parts := strings.Split(callback.Data, "_")
		if len(parts) != 3 {
			return
		}
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
	if strings.HasPrefix(callback.Data, "more_") {
		parts := strings.Split(callback.Data, "_")
		if len(parts) != 2 {
			return
		}
		referralCode := parts[1]

		details, err := getMoreDetails(referralCode)
		if err != nil {
			log.Printf("Error retrieving more details: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Xatolik yuz berdi. Iltimos, keyinroq yana urinib ko'ring.\nAn error occurred. Please try again later.")
			bot.Send(msg)
			return
		}

		detailMsg := tgbotapi.NewMessage(callback.Message.Chat.ID, details)
		detailMsg.ParseMode = "HTML"
		if _, err := bot.Send(detailMsg); err != nil {
			log.Printf("Error sending detailed message: %v", err)
		}
	}
}

func handleBonusForExamify(message *tgbotapi.Message) {
	referralCode := fmt.Sprintf("REF%d", message.From.ID)

	err := storeCoupon(message.Chat.ID, referralCode, message.Chat.FirstName)
	if err != nil {
		log.Printf("Failed to store coupon: %v", err)
		return
	}

	botUsername := bot.Self.UserName
	startLink := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, referralCode)

	responseMessage := fmt.Sprintf(`
🇺🇿 Sizga examify.uz dan maxsus bonus kaliti taqdim etildi! 🤩 Ushbu havola orqali kirgan har 2 ta foydalanuvchi uchun siz examify.uz dan 1 ta sifatli to'liq mock imtihoniga ega bo'lasiz. Bu imkoniyatni qo'ldan boy bermang va do'stlaringizni ham qo'shiling! 🎯

📌 Sizning referal havolangiz: %s

Do'stlaringiz bilan ulashing va imtihonlarga tayyorgarlik ko'rishni qiziqarli va samarali qiling!

🇺🇸 You have received a special bonus key for examify.uz! 🎉 For every 2 people who sign up through your link, you'll get access to a high-quality full mock exam on examify.uz. Don't miss out on this chance and invite your friends to join in! 📚

🔗 Your referral link: %s

Share with your friends and make exam preparation fun and effective!
`, startLink, startLink)

	msg := tgbotapi.NewMessage(message.Chat.ID, responseMessage)

	shareButton := tgbotapi.NewInlineKeyboardButtonSwitch("🔗 Share", responseMessage)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(shareButton),
	)
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func handleDetailsOfBonus(msg *tgbotapi.Message) {
	responseMsg, keyboard := getDetailsOfBonus(msg.Chat.ID)
	messageConfig := tgbotapi.NewMessage(msg.Chat.ID, responseMsg)
	messageConfig.ParseMode = "HTML"
	if keyboard != nil {
		messageConfig.ReplyMarkup = keyboard
	}
	_, err := bot.Send(messageConfig)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
