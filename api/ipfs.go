package api

import (
	"github.com/gin-gonic/gin"
	service "spike-blockchain-server/service/ipfs"
)

func PinFile(c *gin.Context, filePath string) {
	var service service.PinFileService
	service.FilePath = filePath
	if err := c.ShouldBind(&service); err == nil {
		res := service.PinFile()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}

func PinJSON(c *gin.Context, Json string) {
	var service service.PinJsonService
	service.Json = Json
	if err := c.ShouldBind(&service); err == nil {
		res := service.PinJson()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}

func DownLoadFile(c *gin.Context, CID string, savePath string) {
	var service service.DownloadFileService
	service.CID = CID
	service.SaveFilePath = savePath
	if err := c.ShouldBind(&service); err == nil {
		res := service.DownloadFile()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}
