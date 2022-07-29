package server

import (
	"github.com/gin-gonic/gin"

	"spike-blockchain-server/api"
	"spike-blockchain-server/chain"
)

func NewRouter(chainApi *chain.BscListener) *gin.Engine {
	r := gin.Default()

	// to use
	//r.Use(middleware.EthSignatureVerify())
	//r.Use(middleware.ApiKeyAuth())
	//r.Use(middleware.LoggerToFile())

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
			ipfs.POST("pin/jsonfile", api.PinJSONFile)
			ipfs.GET("file", api.DownloadFile)
			ipfs.GET("json", api.DownloadJSON)
		}
		chain := v1.Group("/chain")
		{
			chain.POST("tx/isPending", chainApi.QueryTxIsPendingByHash)
			chain.POST("tx/status", chainApi.QueryTxStatusByHash)
			chain.POST("nft/metadata", chainApi.QueryNftMetadata)
		}
		wallet := v1.Group("/wallet")
		{
			wallet.GET("erc20", api.ERC20TxRecord)
			wallet.GET("native", api.NativeTxRecord)
		}
	}
	return r
}
