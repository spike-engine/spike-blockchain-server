package nft

import (
	"github.com/joho/godotenv"
	"testing"
)

func init() {
	godotenv.Load("../../.env")
}

func TestEthAllNFTsService(t *testing.T) {
	service := FindAllNFTsService{
		Chain:   "eth",
		Format:  "decimal",
		Limit:   50,
		Address: "0x4f025c68c27f860946dc9fa2814bef2a6a7e2ce0",
	}
	service.FindAllNFTs()
}

func TestEthAllNFTsServiceWithSmallAmount(t *testing.T) {
	service := FindAllNFTsService{
		Chain:   "eth",
		Format:  "decimal",
		Address: "0x3FE63F0f8469497223A9eb800be62D75b8B8e6Eb",
	}
	service.FindAllNFTs()
}
