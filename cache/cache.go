package cache

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(
		&redis.Options{
			Addr:       os.Getenv("REDIS_ADDR"),
			Password:   os.Getenv("REDIS_PW"),
			DB:         int(db),
			MaxRetries: 1,
		})

	_, err := client.Ping().Result()

	if err != nil {
		// log connecting to redis failed
	}

	RedisClient = client
}
