package cache

import (
	"github.com/go-redis/redis"
	logger "github.com/ipfs/go-log"
	"os"
	"strconv"
)

var log = logger.Logger("cache")

var RedisClient *redis.Client

func Redis() error {
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
		log.Error("redis init err : ", err)
		return err
	}

	RedisClient = client
	return nil
}
