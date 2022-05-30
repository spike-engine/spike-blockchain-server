package nft

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"

	nftserializer "spike-blockchain-server/serializer/nft"
	"spike-blockchain-server/service"
)

type FindAllNFTsService struct {
	Chain   string `json:"chain"`
	Format  string `json:"format,omitempty"`
	Limit   int32  `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	Address string `json:"address"`
}

func (s *FindAllNFTsService) FindAllNFTs() {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("x-api-key", os.Getenv("MORALIS_KEY")).
		Get(s.url())
	if err != nil {
		log.Panic(err)
	}

	var res nftserializer.MoralisNFTsResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		log.Panic(err)
	}
}

func (s *FindAllNFTsService) url() string {
	return fmt.Sprintf("%s%s/nft?chain=%s&format=%s&limit=%d&cursor=%s", service.MORALIS_API, s.Address, s.Chain, s.Format, s.Limit, s.Cursor)
}
