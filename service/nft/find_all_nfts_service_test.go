package nft

import (
	"testing"
)

func TestEthAllNftsService(t *testing.T) {
	service := FindAllNftsService{
		Chain:   "eth",
		Format:  "decimal",
		Address: "0x3FE63F0f8469497223A9eb800be62D75b8B8e6Eb",
	}
	service.FindAllEthNfts()
}
