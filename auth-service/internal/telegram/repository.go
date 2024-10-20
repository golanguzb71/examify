package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"time"
)

type Coupon struct {
	ChatID       int64  `bson:"chat_id"`
	ReferralCode string `bson:"referral_code"`
	WelcomeCount int    `bson:"welcome_count"`
	UsedCount    int    `bson:"used_count"`
	Name         string `bson:"name"`
	CreatedAt    string `bson:"created_at"`
}

type LinkedUserCollection struct {
	ChatId             int64  `bson:"chat_id"`
	LinkedReferralCode string `bson:"linked_referral_code"`
	Name               string `bson:"name"`
	CreatedAt          string `bson:"created_at"`
}

func storeCoupon(chatID int64, referralCode, name string) error {
	couponCollection := bonusdb.Collection("couponCollection")

	var existingUserCoupon Coupon
	err := couponCollection.FindOne(context.TODO(), bson.M{"chat_id": chatID}).Decode(&existingUserCoupon)
	if err == nil {
		return nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error checking user coupon: %v", err)
	}

	newCoupon := Coupon{
		ChatID:       chatID,
		ReferralCode: referralCode,
		WelcomeCount: 0,
		Name:         name,
		UsedCount:    0,
		CreatedAt:    time.Now().String(),
	}
	_, err = couponCollection.InsertOne(context.TODO(), newCoupon)
	if err != nil {
		return fmt.Errorf("error inserting new coupon: %v", err)
	}

	log.Printf("New coupon created for chatID %d with referral code %s", chatID, referralCode)
	return nil
}

func storeUsedUser(chatId int64, linkedReferralCode string, name string) error {
	storedUserCollection := bonusdb.Collection("linkedUserCollection")
	newLinkedUser := LinkedUserCollection{
		ChatId:             chatId,
		LinkedReferralCode: linkedReferralCode,
		Name:               name,
		CreatedAt:          time.Now().String(),
	}
	_, err := storedUserCollection.InsertOne(context.TODO(), newLinkedUser)
	if err != nil {
		return err
	}
	return nil
}

func processReferral(newUserChatID int64, referralCode, firstName string) error {
	couponCollection := bonusdb.Collection("couponCollection")

	var existingUserCoupon Coupon
	err := couponCollection.FindOne(context.TODO(), bson.M{"chat_id": newUserChatID}).Decode(&existingUserCoupon)
	if err == nil {
		log.Printf("User with chatID %d already exists, ignoring referral", newUserChatID)
		return nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error checking new user: %v", err)
	}

	var referrerCoupon Coupon
	err = couponCollection.FindOne(context.TODO(), bson.M{"referral_code": referralCode}).Decode(&referrerCoupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid referral code: %s", referralCode)
		}
		return fmt.Errorf("error checking referral code: %v", err)
	}

	if referrerCoupon.ChatID == newUserChatID {
		return fmt.Errorf("cannot use your own referral code")
	}

	_, err = couponCollection.UpdateOne(
		context.TODO(),
		bson.M{"referral_code": referralCode},
		bson.M{"$inc": bson.M{"welcome_count": 1}},
	)
	if err != nil {
		return fmt.Errorf("error incrementing welcome count: %v", err)
	}
	err = storeUsedUser(newUserChatID, referralCode, firstName)
	if err != nil {
		return fmt.Errorf("error incrementing welcome count: %v", err)
	}
	log.Printf("Incremented welcome count for referral code %s", referralCode)
	return nil
}

func getDetailsOfBonus(chatId int64) (string, *tgbotapi.InlineKeyboardMarkup) {
	couponCollection := bonusdb.Collection("couponCollection")
	var existingUserCoupon Coupon
	err := couponCollection.FindOne(context.TODO(), bson.M{"chat_id": chatId}).Decode(&existingUserCoupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "Siz hali hech qanday havola ulashishni boshlamadingiz.\nYou haven't started sharing links yet.", nil
		}
		log.Println("Error finding coupon:", err)
		return "Kupon ma'lumotlarini olishda xatolik yuz berdi.\nAn error occurred while retrieving your coupon details.", nil
	}

	createdAtTrimmed := existingUserCoupon.CreatedAt
	if spaceIndex := strings.Index(createdAtTrimmed, " "); spaceIndex != -1 {
		createdAtTrimmed = createdAtTrimmed[:spaceIndex+9]
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtTrimmed)
	if err != nil {
		log.Println("Error parsing CreatedAt time:", err)
		return "Vaqtni formatlashda xatolik yuz berdi.\nAn error occurred while formatting the time.", nil
	}

	formattedTime := createdAt.Format("15:04:05")

	msg := fmt.Sprintf(
		`
🇺🇿
<b>Foydalanuvchi nomingiz:</b>  <code>%s</code>
<b>Sizning referal havolangiz:</b> <code>%s</code>
<b>Siz taklif qilgan umumiy foydalanuvchilar soni:</b> <code>%d</code>
<b>Siz ishlatib yuborgan kuponlaringiz soni:</b> <code>%d</code>
<b>Ro'yxatdan o'tgan vaqtingiz:</b> <code>%s</code>

🏴󠁧󠁢󠁥󠁮󠁧󠁿
<b>Your Username:</b> <code>%s</code>
<b>Your Referral Link:</b> <code>%s</code>
<b>Total Users You Invited:</b> <code>%d</code>
<b>Coupons You Have Used:</b> <code>%d</code>
<b>Registration Time:</b> <code>%s</code>`,
		existingUserCoupon.Name,
		existingUserCoupon.ReferralCode,
		existingUserCoupon.WelcomeCount,
		existingUserCoupon.UsedCount,
		formattedTime,
		existingUserCoupon.Name,
		existingUserCoupon.ReferralCode,
		existingUserCoupon.WelcomeCount,
		existingUserCoupon.UsedCount,
		formattedTime,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ko'proq (More Info)", fmt.Sprintf("more_%s", existingUserCoupon.ReferralCode)),
		),
	)

	return msg, &keyboard
}

func getMoreDetails(referralCode string) (string, error) {
	linkedUserCollection := bonusdb.Collection("linkedUserCollection")
	var linkedUsers []LinkedUserCollection

	filter := bson.M{"linked_referral_code": referralCode}
	cursor, err := linkedUserCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding linked users: %v", err)
		return "Xatolik yuz berdi. Iltimos, keyinroq yana urinib ko'ring.\nAn error occurred. Please try again later.", err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var linkedUser LinkedUserCollection
		err := cursor.Decode(&linkedUser)
		if err != nil {
			log.Printf("Error decoding linked user: %v", err)
			continue
		}
		linkedUsers = append(linkedUsers, linkedUser)
	}

	if len(linkedUsers) == 0 {
		return "Hech qanday siz taklif qilgan foydalanuvchilar topilmadi.\nNo users found who you invited.", nil
	}

	var msgBuilder strings.Builder
	msgBuilder.WriteString("🇺🇿\n<b>Ulanish havolasi orqali qo'shilgan foydalanuvchilar:</b>\n")
	for i, user := range linkedUsers {
		msgBuilder.WriteString(fmt.Sprintf(
			"%d)  Foydalanuvchi nomi: <code>%s</code>, Ro'yxatdan o'tgan vaqti: <code>%s</code>\n",
			i+1,
			user.Name,
			user.CreatedAt[:19],
		))
	}

	msgBuilder.WriteString("\n🏴󠁧󠁢󠁥󠁮󠁧󠁿\n<b>Users joined via your referral link:</b>\n")
	for i, user := range linkedUsers {
		msgBuilder.WriteString(fmt.Sprintf(
			"%d) User name: <code>%s</code>, Registration time: <code>%s</code>\n",
			i+1,
			user.Name,
			user.CreatedAt[:19],
		))
	}

	return msgBuilder.String(), nil
}
