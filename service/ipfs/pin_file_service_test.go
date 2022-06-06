package ipfs

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func TestPinFileDuplicate(t *testing.T) {
	service := PinFileService{
		FilePath: "../../.env.example",
	}

	res := service.PinFile()
	fmt.Println(res)
}
