package telegram

import (
	database "authService/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"time"
)

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func storeCode(code string, chatID int64) {
	database.RDB.Set(context.TODO(), code, chatID, 1*time.Minute)
	database.RDB.Set(context.TODO(), fmt.Sprintf("%v", chatID), code, 1*time.Minute)
	log.Printf("Storing code %s for chat ID %d", code, chatID)
}

func GetStoredCode(key string) *string {
	result, err := database.RDB.Get(context.TODO(), fmt.Sprintf("%v", key)).Result()
	fmt.Println(result)
	if errors.Is(err, redis.Nil) {
		log.Printf("No code stored for key: %v", key)
		return nil
	} else if err != nil {
		log.Printf("Error fetching code for key: %v, error: %v", key, err)
		return nil
	}
	return &result
}
