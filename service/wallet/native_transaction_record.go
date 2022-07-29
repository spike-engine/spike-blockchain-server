package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	"os"
	"spike-blockchain-server/constants"
	"spike-blockchain-server/serializer"
)

type NativeTransactionRecordService struct {
	Address string `form:"address" json:"address" binding:"required"`
}

func (s *NativeTransactionRecordService) FindNativeTransactionRecord() serializer.Response {
	client := resty.New()

	apiKey := os.Getenv("BSC_API_KEY")

	bscClient, err := ethclient.Dial(os.Getenv("MORALIS_SPEEDY_NODE"))
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}
	defer bscClient.Close()
	blockNumber, err := bscClient.BlockNumber(context.Background())
	if err != nil {
		return serializer.Response{
			Code:  402,
			Error: err.Error(),
		}
	}

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(s.url(apiKey, blockNumber))
	if err != nil {
		return serializer.Response{
			Code:  403,
			Error: err.Error(),
		}
	}
	var bscRes BscRes
	var bnbRecord []Result

	err = json.Unmarshal(resp.Body(), &bscRes)
	if err != nil {
		return serializer.Response{
			Code:  405,
			Error: err.Error(),
		}
	}
	if len(bscRes.Result) != 0 {

		for i := range bscRes.Result {
			if bscRes.Result[i].Input == "0x" {
				bnbRecord = append(bnbRecord, bscRes.Result[i])
			}
		}
	}

	bscRes.Result = bnbRecord
	return serializer.Response{
		Code: 200,
		Data: bscRes,
	}
}

func (s *NativeTransactionRecordService) url(apiKey string, blockNumber uint64) string {
	return fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&offset=10000&page=1&sort=desc&apikey=%s", constants.BSCSCAN_API, s.Address, blockNumber-201600, blockNumber, apiKey)
}
