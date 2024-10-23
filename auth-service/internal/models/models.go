package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ChatID       int64              `bson:"chat_id" json:"chat_id"`
	CreatedAt    string             `bson:"created_at" json:"created_at"`
	Name         string             `bson:"name" json:"name"`
	ReferralCode string             `bson:"referral_code" json:"referral_code"`
	UsedCount    int32              `bson:"used_count" json:"used_count"`
	WelcomeCount int32              `bson:"welcome_count" json:"welcome_count"`
}

type LinkedUser struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ChatID             int64              `bson:"chat_id" json:"chat_id"`
	CreatedAt          string             `bson:"created_at" json:"created_at"`
	LinkedReferralCode string             `bson:"linked_referral_code" json:"linked_referral_code"`
	Name               string             `bson:"name" json:"name"`
}

type GetMoreBonusInformation struct {
	GuestName         string `json:"guestName" protobuf:"bytes,1,opt,name=guestName,proto3"`
	GuestChatID       string `json:"guestChatId" protobuf:"bytes,2,opt,name=guestChatId,proto3"`
	GuestRegisteredAt string `json:"guestRegisteredAt" protobuf:"bytes,3,opt,name=guestRegisteredAt,proto3"`
}

type ReferralInfo struct {
	Name         string                     `json:"name" protobuf:"bytes,1,opt,name=name,proto3"`
	RefLink      string                     `json:"refLink" protobuf:"bytes,2,opt,name=refLink,proto3"`
	WelcomeCount int32                      `json:"welcomeCount" protobuf:"varint,3,opt,name=welcomeCount,proto3"`
	BonusCount   int32                      `json:"bonusCount" protobuf:"varint,4,opt,name=bonusCount,proto3"`
	RegisteredAt string                     `json:"registeredAt" protobuf:"bytes,5,opt,name=registeredAt,proto3"`
	More         []*GetMoreBonusInformation `json:"more" protobuf:"bytes,6,rep,name=more,proto3"`
}

const (
	CouponsCollection     = "couponCollection"
	LinkedUsersCollection = "linkedUserCollection"
)
