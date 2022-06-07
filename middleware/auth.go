package middleware

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"

	"spike-blockchain-server/serializer"
)

func EthSignatureVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params ethParams
		if err := c.ShouldBind(&params); err == nil {
			if params.verify() {
				c.Next()
			} else {
				c.JSON(200, serializer.Response{
					Code:  101,
					Error: "Invalid signature",
				})
			}
		} else {
			c.JSON(200, serializer.Response{
				Code:  100,
				Error: err.Error(),
			})
			c.Abort()
		}
	}
}

type ethParams struct {
	Chain string `form:"chain" binding:"required"`
	Sign  string `form:"sign"  binding:"required"`
	Nonce string `form:"nonce" binding:"required"`
	Addr  string `form:"addr"  binding:"required"`
}

func (params *ethParams) verify() bool {
	sig := hexutil.MustDecode(params.Sign)
	msg := accounts.TextHash([]byte(params.Nonce))
	sig[crypto.RecoveryIDOffset] -= 27

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return strings.ToLower(params.Addr) == strings.ToLower(recoveredAddr.String())
}
