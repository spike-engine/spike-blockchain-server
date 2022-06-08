package ipfs

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../.env")
}

func TestPinFileDuplicate(t *testing.T) {
	service := PinFileService{
		Filepath: "../../.env.example",
	}

	assert.Equal(t, service.PinFile().Code, 200)
}
