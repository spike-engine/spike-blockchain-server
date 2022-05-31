package nft

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/serializer"
	"spike-blockchain-server/service"
)

type FindAllNFTsService struct {
	Chain   string `form:"chain"   json:"chain"   binding:"required"`
	Format  string `form:"format"  json:"format,omitempty"`
	Limit   int32  `form:"limit"   json:"limit,omitempty"`
	Cursor  string `form:"cursor"  json:"cursor,omitempty"`
	Address string `form:"address" json:"address" binding:"required"`
}

func (s *FindAllNFTsService) FindAllNFTs() serializer.Response {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("x-api-key", os.Getenv("MORALIS_KEY")).
		Get(s.url())
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return serializer.Response{
			Code:  resp.StatusCode(),
			Error: resp.String(),
		}
	}

	var res serializer.MoralisNFTs
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return serializer.Response{
			Code:  401,
			Error: err.Error(),
		}
	}

	if resp.IsSuccess() {
		return serializer.Response{
			Code: resp.StatusCode(),
			Data: res,
		}
	}

	return serializer.Response{
		Code: 300,
		Data: resp.String(),
	}
}

func (s *FindAllNFTsService) url() string {
	return fmt.Sprintf("%s%s/nft?chain=%s&format=%s&limit=%d&cursor=%s", service.MORALIS_API, s.Address, s.Chain, s.Format, s.Limit, s.Cursor)
}
