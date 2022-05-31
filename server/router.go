package server

import (
	"github.com/gin-gonic/gin"
	"spike-blockchain-server/api"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})

		nft := v1.Group("/nft")
		{
			nft.GET("all", api.FindAllNFTs)
		}
	}
	return r
}
