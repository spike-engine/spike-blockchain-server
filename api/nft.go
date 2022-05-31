package api

import (
	"github.com/gin-gonic/gin"
	service "spike-blockchain-server/service/nft"
)

func FindAllNFTs(c *gin.Context) {
	var service service.FindAllNFTsService
	if err := c.ShouldBind(&service); err == nil {
		res := service.FindAllNFTs()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}
