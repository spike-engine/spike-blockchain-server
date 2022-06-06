package ipfs

import (
	"fmt"
	"github.com/joho/godotenv"
	"testing"
)

func init() {
	godotenv.Load("../../.env")
}

func TestDownLoadFile(t *testing.T) {
	service := DownloadFileService{
		CID:          "QmfNxKabtsdNgdLCib1AXaVrMvyLaFD2UbMmsmoRNJTU8w",
		SaveFilePath: "/Users/fuyiwei/Desktop/img.txt",
	}

	res := service.DownloadFile()
	fmt.Printf("%+v", res)
}
