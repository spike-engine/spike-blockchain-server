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
	var service service.PinJsonService
	if err := c.ShouldBind(&service); err == nil {
		res := service.PinJson()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}

func DownloadFile(c *gin.Context) {
	var service service.DownloadService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Download()
		if res.Code == 200 {
			c.Writer.Write(res.Data.([]byte))
		} else {
			c.JSON(200, res)
		}
	} else {
		c.JSON(200, err.Error())
	}
}

func DownloadJSON(c *gin.Context) {
	var service service.DownloadService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Download()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}
