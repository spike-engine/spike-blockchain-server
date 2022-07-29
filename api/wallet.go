package api

import (
	"github.com/gin-gonic/gin"
	service "spike-blockchain-server/service/wallet"
)

func ERC20TxRecord(c *gin.Context) {
	var service service.ERC20TransactionRecordService
	if err := c.ShouldBind(&service); err == nil {
		res := service.FindERC20TransactionRecord()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}

func NativeTxRecord(c *gin.Context) {
	var service service.NativeTransactionRecordService
	if err := c.ShouldBind(&service); err == nil {
		res := service.FindNativeTransactionRecord()
		c.JSON(200, res)
	} else {
		c.JSON(200, err.Error())
	}
}
