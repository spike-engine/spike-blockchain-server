package nft

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	serialize "spike-blockchain-server/serialize/nft"
	"spike-blockchain-server/service"
)

type FindAllNftsService struct {
	Chain   string `json:"chain"`
	Format  string `json:"format,omitempty"`
	Limit   int32  `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	Address string `json:"address"`
}

func (s *FindAllNftsService) FindAllEthNfts() {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("x-api-key", "7L9fFRFaUpKOlFoUZi7SN0UBiUWbHgf6tTenyqg19TjPRTGuRS46UMRierHQ3XIs").
		Get(s.url())
	if err != nil {
		log.Panic(err)
	}

	var res serialize.FindAllNftsResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		log.Panic(err)
	}
	res.Total = len(res.Result)
	nfts := res.Result

	_ = nfts
}

func (s *FindAllNftsService) url() string {
	return fmt.Sprintf("%s%s/nft?chain=%s&format=%s&limit=%d&cursor=%s", service.MORALIS_API, s.Address, s.Chain, s.Format, s.Limit, s.Cursor)
}
