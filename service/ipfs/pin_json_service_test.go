package ipfs

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func TestPinJsonDuplicate(t *testing.T) {
	service := PinJsonService{
		Json: `{"pinata":["file","json"],"len":2}`,
	}

	res := service.PinJson()
	fmt.Println(res)
}
