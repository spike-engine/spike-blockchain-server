package api

import (
	"github.com/gin-gonic/gin"
	service "spike-blockchain-server/service/ipfs"
)

func PinFile(c *gin.Context) {
	var service service.PinFileService
	if err := c.ShouldBind(&service); err == nil {
		res := service.PinFile()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}

func PinJSON(c *gin.Context) {

}
