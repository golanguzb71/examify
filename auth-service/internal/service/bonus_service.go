package service

import (
	"authService/proto/pb"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
