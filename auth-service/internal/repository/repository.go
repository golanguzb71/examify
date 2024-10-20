package database

import (
	"authService/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var RDB *redis.Client

func ConnectRedis(cfg *config.Config) {
	redisHost := cfg.Redis.Host
	redisPort := cfg.Redis.Port

	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	fmt.Println("Redis connection established successfully")
}

func ConnectMongo(databasePort, databaseHost, databaseName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s/%s", databaseHost, databasePort, databaseName)
	clientOptions := options.Client().ApplyURI(uri)

	mongodbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = mongodbClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	database := mongodbClient.Database(databaseName)

	log.Printf("Connected to MongoDB database '%s' successfully", databaseName)
	return database, nil
}
