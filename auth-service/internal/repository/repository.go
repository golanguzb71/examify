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

func ConnectMongo(databasePort, databaseHost, databaseName, databasePassword, databaseUsername string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the MongoDB URI with authentication credentials
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", databaseUsername, databasePassword, databaseHost, databasePort, databaseName)

	clientOptions := options.Client().ApplyURI(uri)

	// Connect to the MongoDB client
	mongodbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to verify the connection
	err = mongodbClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Check if the database client is nil before using it
	if mongodbClient == nil {
		return nil, fmt.Errorf("MongoDB client is nil")
	}

	// Access the specified database
	database := mongodbClient.Database(databaseName)
	if database == nil {
		return nil, fmt.Errorf("MongoDB database is nil")
	}

	log.Printf("Connected to MongoDB database '%s' successfully", databaseName)
	return database, nil
}
