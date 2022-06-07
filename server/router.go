package server

import (
	"github.com/gin-gonic/gin"
	"spike-blockchain-server/api"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// to use
	//r.Use(middleware.EthSignatureVerify())

	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})

		nft := v1.Group("/nft")
		{
			nft.GET("all", api.FindAllNFTs)
		}

		ipfs := v1.Group("/ipfs")
		{
			ipfs.POST("pin/file", api.PinFile)
			ipfs.POST("pin/json", api.PinJSON)
			ipfs.GET("file", api.DownloadFile)
			ipfs.GET("json", api.DownloadJSON)
		}
	}
	return r
}
