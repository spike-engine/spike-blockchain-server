package nft

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"spike-blockchain-server/serializer"
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
	assert.NotEmpty(t, service.FindAllNFTs().Data.(serializer.MoralisNFTs).Total)
}

func TestEthAllNFTsServiceWithSmallAmount(t *testing.T) {
	service := FindAllNFTsService{
		Chain:   "eth",
		Format:  "decimal",
		Address: "0x3FE63F0f8469497223A9eb800be62D75b8B8e6Eb",
	}
	assert.NotEmpty(t, service.FindAllNFTs().Data.(serializer.MoralisNFTs).Total)
}

func TestNoETHNFT(t *testing.T) {
	service := FindAllNFTsService{
		Chain:   "eth",
		Format:  "decimal",
		Address: "0x33AD388F713f7A043504f9c0b717841aC5a34e0B",
	}
	res := service.FindAllNFTs()
	assert.Equal(t, res.Code, 200)
	assert.Equal(t, res.Data, serializer.MoralisNFTs{Total: 0, Result: []serializer.NFT{}})
}
