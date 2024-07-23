package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient       *redis.Client
	defaultIPLimit    int
	defaultTokenLimit int
	blockDuration     time.Duration
}

func NewRateLimiter() (*RateLimiter, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	defaultIPLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_IP_LIMIT"))
	defaultTokenLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_TOKEN_LIMIT"))
	blockDuration, _ := time.ParseDuration(os.Getenv("BLOCK_DURATION"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &RateLimiter{
		redisClient:       rdb,
		defaultIPLimit:    defaultIPLimit,
		defaultTokenLimit: defaultTokenLimit,
		blockDuration:     blockDuration,
	}, nil
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, key string, limit int) (bool, error) {
	val, err := rl.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Log key value after increment
	fmt.Printf("Key: %s, Value after Incr: %d, Limit: %d\n", key, val, limit)

	if val == 1 {
		err = rl.redisClient.Expire(ctx, key, time.Second).Err()
		if err != nil {
			return false, err
		}
		fmt.Printf("Set expiration for key: %s, duration: %s\n", key, time.Second)
	}

	if val > int64(limit) {
		err = rl.redisClient.Set(ctx, key, val, rl.blockDuration).Err()
		if err != nil {
			return false, err
		}
		fmt.Printf("Exceeded limit for key: %s, extended expiration to: %s\n", key, rl.blockDuration)
		return false, nil
	}

	return true, nil
}

func (rl *RateLimiter) AllowIP(ctx context.Context, ip string) (bool, error) {
	return rl.AllowRequest(ctx, "ip:"+ip, rl.defaultIPLimit)
}

func (rl *RateLimiter) AllowToken(ctx context.Context, token string) (bool, error) {
	return rl.AllowRequest(ctx, "token:"+token, rl.defaultTokenLimit)
}

func (rl *RateLimiter) GetDefaultIPLimit() int {
	return rl.defaultIPLimit
}

func (rl *RateLimiter) GetDefaultTokenLimit() int {
	return rl.defaultTokenLimit
}
