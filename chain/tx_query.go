package chain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"golang.org/x/xerrors"
	"math/big"
	"spike-blockchain-server/chain/contract"
	"spike-blockchain-server/serializer"
	"strconv"
	"strings"
)

type txQueryService struct {
	TxHash string `json:"txHash"`
}

type metadataService struct {
	TokenId string `json:"tokenId"`
	Address string `json:"address"`
}

func (bl *BscListener) QueryTxIsPendingByHash(c *gin.Context) {
	var service txQueryService
	if err := c.ShouldBind(&service); err == nil {
		log.Infof("tx: %s", service.TxHash)
		res := bl.queryTxIsPendingByHash(service.TxHash)
		c.JSON(200, res)
	} else {
		c.JSON(500, err.Error())
	}
}

func (bl *BscListener) QueryTxStatusByHash(c *gin.Context) {
	var service txQueryService
	if err := c.ShouldBind(&service); err == nil {
		log.Infof("tx: %s", service.TxHash)
		res := bl.queryTxStatusByHash(service.TxHash)
		c.JSON(200, res)
	} else {
		c.JSON(500, err.Error())
	}
}

func (bl *BscListener) queryTxIsPendingByHash(txHash string) serializer.Response {
	_, isPending, err := bl.ec.TransactionByHash(context.Background(), common.HexToHash(txHash))
	code := 200
	if err != nil {
		code = 500
		return serializer.Response{
			Code:  code,
			Data:  isPending,
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: code,
		Data: isPending,
	}
}

func (bl *BscListener) queryTxStatusByHash(txHash string) serializer.Response {
	receipt, err := bl.ec.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	code := 200
	if err != nil {
		code = 500
		return serializer.Response{
			Code:  code,
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: code,
		Data: receipt.Status,
	}
}

func (bl *BscListener) QueryNftMetadata(c *gin.Context) {
	var service metadataService

	if err := c.ShouldBind(&service); err == nil {
		tokenId, err := strconv.Atoi(service.TokenId)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		res := bl.queryNftMetadata(int64(tokenId), service.Address)
		c.JSON(200, res)
	} else {
		c.JSON(500, err.Error())
	}
}

func (bl *BscListener) queryNftMetadata(tokenId int64, address string) serializer.Response {
	aunft, err := contract.NewAunft(common.HexToAddress(AUNFTContractAddress), bl.ec)
	if err != nil {
		fmt.Println(err)
		return serializer.Response{
			Code:  500,
			Error: err.Error(),
		}
	}
	uri, err := aunft.TokenURI(nil, big.NewInt(tokenId))
	if err != nil {
		fmt.Println(err)
		return serializer.Response{
			Code:  500,
			Error: err.Error(),
		}
	}
	owner, err := aunft.OwnerOf(nil, big.NewInt(tokenId))
	if err != nil {
		return serializer.Response{
			Code:  500,
			Error: err.Error(),
		}
	}
	if strings.ToLower(owner.String()) != strings.ToLower(address) {
		return serializer.Response{
			Code:  500,
			Error: xerrors.New("tokenId, nft not match").Error(),
		}
	}

	client := resty.New()

	resp, err := client.R().Get(uri)
	log.Infof("query nft tokenId : %d, uri : %s", tokenId, uri)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Error: err.Error(),
		}
	}

	if resp.IsError() {
		return serializer.Response{
			Code:  500,
			Error: resp.String(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: string(resp.Body()),
	}
}
