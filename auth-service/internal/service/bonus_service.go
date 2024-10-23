package service

import (
	"authService/proto/pb"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strconv"
)

type BonusService struct {
	pb.UnimplementedBonusServiceServer
	db *mongo.Database
}

func NewBonusService(db *mongo.Database) *BonusService {
	return &BonusService{db: db}
}
func (b *BonusService) UseBonusAttempt(ctx context.Context, req *pb.UseBonusAttemptRequest) (*pb.UseBonusAttemptResponse, error) {
	chatId := req.ChatId
	couponCollection := b.db.Collection("couponCollection")
	var coupon struct {
		UsedCount    int `bson:"used_count"`
		WelcomeCount int `bson:"welcome_count"`
	}

	err := couponCollection.FindOne(ctx, bson.M{"chat_id": chatId}).Decode(&coupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &pb.UseBonusAttemptResponse{
				Response: false,
			}, nil
		}
		return nil, fmt.Errorf("error fetching coupon details: %v", err)
	}

	usedCount := coupon.UsedCount
	welcomeCount := coupon.WelcomeCount
	var resp bool
	if usedCount < welcomeCount/2 {
		resp = true
		_, _ = couponCollection.UpdateOne(
			ctx,
			bson.M{"chat_id": chatId},
			bson.M{"$inc": bson.M{"used_count": 1}},
		)
	} else {
		resp = false
	}
	return &pb.UseBonusAttemptResponse{
		Response: resp,
	}, nil
}

func (b *BonusService) CalculateBonusByChatId(ctx context.Context, req *pb.BonusServiceAbsRequest) (*pb.CalculateBonusByChatIdResponse, error) {
	chatId := req.ChatId
	couponCollection := b.db.Collection("couponCollection")

	chatIdFloat64, err := strconv.ParseFloat(chatId, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chatId format, expected float64: %v", err)
	}

	chatIdInt64 := int64(math.Floor(chatIdFloat64))

	var coupon struct {
		UsedCount    int `bson:"used_count"`
		WelcomeCount int `bson:"welcome_count"`
	}

	fmt.Println(chatIdInt64)

	err = couponCollection.FindOne(ctx, bson.M{"chat_id": chatIdInt64}).Decode(&coupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &pb.CalculateBonusByChatIdResponse{
				Count: 0,
			}, nil
		}
		return nil, fmt.Errorf("error fetching coupon details: %v", err)
	}

	usedCount := coupon.UsedCount
	welcomeCount := coupon.WelcomeCount

	return &pb.CalculateBonusByChatIdResponse{Count: int32(welcomeCount/2 - usedCount)}, nil
}

func (b *BonusService) GetBonusInformationByChatId(ctx context.Context, req *pb.BonusServiceAbsRequest) (*pb.GetBonusInformationByChatIdResponse, error) {
	chatId := req.ChatId
	couponCollection := b.db.Collection("couponCollection")
	linkedUserCollection := b.db.Collection("linkedUserCollection")
	chatIdFloat64, err := strconv.ParseFloat(chatId, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chatId format, expected float64: %v", err)
	}
	chatIdInt64 := int64(math.Floor(chatIdFloat64))
	var coupon struct {
		Name               string `bson:"name"`
		WelcomeCount       int    `bson:"welcome_count"`
		UsedCount          int    `bson:"used_count"`
		CreatedAt          string `bson:"created_at"`
		LinkedReferralCode string `bson:"linked_referral_code"`
	}
	err = couponCollection.FindOne(ctx, bson.M{"chat_id": chatIdInt64}).Decode(&coupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no bonus information found for chat_id: %v", chatIdInt64)
		}
		return nil, fmt.Errorf("error fetching coupon details: %v", err)
	}

	response := &pb.GetBonusInformationByChatIdResponse{
		Name:         coupon.Name,
		RefLink:      coupon.LinkedReferralCode,
		WelcomeCount: int32(coupon.WelcomeCount),
		BonusCount:   int32(coupon.WelcomeCount/2 - coupon.UsedCount),
		RegisteredAt: coupon.CreatedAt,
	}
	if coupon.LinkedReferralCode != "" {
		cursor, err := linkedUserCollection.Find(ctx, bson.M{"referral_code": coupon.LinkedReferralCode})
		if err != nil {
			return nil, fmt.Errorf("error fetching linked user information: %v", err)
		}
		defer cursor.Close(ctx)

		var moreBonusInfo []*pb.GetMoreBonusInformation
		for cursor.Next(ctx) {
			var linkedUser struct {
				Name      string `bson:"name"`
				ChatId    int64  `bson:"chat_id"`
				CreatedAt string `bson:"created_at"`
			}
			err := cursor.Decode(&linkedUser)
			if err != nil {
				return nil, fmt.Errorf("error decoding linked user information: %v", err)
			}
			moreBonusInfo = append(moreBonusInfo, &pb.GetMoreBonusInformation{
				GuestName:         linkedUser.Name,
				GuestChatId:       strconv.FormatInt(linkedUser.ChatId, 10),
				GuestRegisteredAt: linkedUser.CreatedAt,
			})
		}
		response.More = moreBonusInfo
	}

	return response, nil
}
