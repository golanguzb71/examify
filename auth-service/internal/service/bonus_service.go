package service

import (
	"authService/internal/models"
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
	chatIdFloat64, err := strconv.ParseFloat(chatId, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chatId format, expected float64: %v", err)
	}

	chatIdInt64 := int64(math.Floor(chatIdFloat64))
	var coupon models.Coupon
	err = b.db.Collection(models.CouponsCollection).FindOne(ctx, bson.M{
		"chat_id": chatIdInt64,
	}).Decode(&coupon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding coupon: %v", err)
	}
	cursor, err := b.db.Collection(models.LinkedUsersCollection).Find(ctx, bson.M{
		"linked_referral_code": coupon.ReferralCode,
	})
	if err != nil {
		return nil, fmt.Errorf("error finding linked users: %v", err)
	}
	defer cursor.Close(ctx)

	var linkedUsers []models.LinkedUser
	if err := cursor.All(ctx, &linkedUsers); err != nil {
		return nil, fmt.Errorf("error decoding linked users: %v", err)
	}
	moreBonusInfo := make([]*pb.GetMoreBonusInformation, 0, len(linkedUsers))
	for _, user := range linkedUsers {
		moreBonusInfo = append(moreBonusInfo, &pb.GetMoreBonusInformation{
			GuestName:         user.Name,
			GuestChatId:       fmt.Sprintf("%d", user.ChatID),
			GuestRegisteredAt: user.CreatedAt,
		})
	}

	response := &pb.GetBonusInformationByChatIdResponse{
		Name:         coupon.Name,
		RefLink:      coupon.ReferralCode,
		WelcomeCount: coupon.WelcomeCount,
		BonusCount:   coupon.UsedCount,
		RegisteredAt: coupon.CreatedAt,
		More:         moreBonusInfo,
	}

	return response, nil
}
