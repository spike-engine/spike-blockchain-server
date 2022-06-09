package ipfs

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"spike-blockchain-server/serializer"
)

func init() {
	godotenv.Load("../../.env")
}

func TestPinJsonDuplicate(t *testing.T) {
	service := PinJsonService{
		Json: `{"pinata":["file","json"],"len":2}`,
	}

	assert.Equal(t, service.PinJson().Code, 200)
	assert.Equal(t, service.PinJson().Data.(serializer.Pin).IsDuplicate, true)
}
