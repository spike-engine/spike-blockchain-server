package middleware

import (
	"os"

	"github.com/gin-gonic/gin"

	"spike-blockchain-server/cache"
	"spike-blockchain-server/serializer"
)

func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params authParams
		if err := c.ShouldBindHeader(&params); err != nil {
			c.JSON(200, serializer.Response{
				Code:  101,
				Error: "header: api_key is required",
			})
			c.Abort()
		} else {
			res, _ := cache.RedisClient.SIsMember(os.Getenv("REDIS_API_KEY_STORE"), params.APIKey).Result()
			if res {
				c.Next()
			} else {
				c.JSON(200, serializer.Response{
					Code:  101,
					Error: "header: api_key doesn't exist",
				})
				c.Abort()
			}
		}
	}
}

type authParams struct {
	APIKey string `header:"api_key" binding:"required"`
}
