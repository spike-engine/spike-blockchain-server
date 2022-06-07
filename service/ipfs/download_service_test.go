package ipfs

import (
	"fmt"
	"github.com/joho/godotenv"
	"testing"
)

func init() {
	godotenv.Load("../../.env")
}

func TestDownloadJSON(t *testing.T) {
	service := DownloadService{
		IpfsHash: "QmfNxKabtsdNgdLCib1AXaVrMvyLaFD2UbMmsmoRNJTU8w",
	}
	fmt.Println(service.Download())
}

func TestDownloadFile(t *testing.T) {
	service := DownloadService{
		IpfsHash: "QmcQUrnWMMGs2rYcQTipSyhwYFNA8Vp65q245mFfQpYXy2",
	}
	fmt.Println(service.Download())
}
