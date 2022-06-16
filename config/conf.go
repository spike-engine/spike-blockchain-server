package config

import (
	"github.com/joho/godotenv"
	
	"spike-blockchain-server/cache"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cache.Redis()
}
