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

type ERC20TransactionRecordService struct {
	Address        string `form:"address" json:"address" binding:"required"`
	ContractAdress string `form:"contract_address" json:"contract_address" binding:"required"`
}

func (s *ERC20TransactionRecordService) FindERC20TransactionRecord() serializer.Response {
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

	err = json.Unmarshal(resp.Body(), &bscRes)
	if err != nil {
		return serializer.Response{
			Code:  405,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: bscRes,
	}
}

func (s *ERC20TransactionRecordService) url(apiKey string, blockNumber uint64) string {
	return fmt.Sprintf("%s?module=account&action=tokentx&address=%s&startblock=%d&endblock=%d&offset=10000&page=1&sort=desc&apikey=%s&contractaddress=%s", constants.BSCSCAN_API, s.Address, blockNumber-201600, blockNumber, apiKey, s.ContractAdress)
}

type Result struct {
	Hash        string `json:"hash"`
	TimeStamp   string `json:"timeStamp"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Input       string `json:"input"`
}

type BscRes struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []Result `json:"result"`
}
