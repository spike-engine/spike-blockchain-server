package ipfs

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../.env")
}

func TestDownloadJSON(t *testing.T) {
	service := DownloadService{
		IpfsHash: "QmfNxKabtsdNgdLCib1AXaVrMvyLaFD2UbMmsmoRNJTU8w",
	}

	assert.Equal(t, service.Download().Code, 200)
}

func TestDownloadFile(t *testing.T) {
	service := DownloadService{
		IpfsHash: "QmcQUrnWMMGs2rYcQTipSyhwYFNA8Vp65q245mFfQpYXy2",
	}

	assert.Equal(t, service.Download().Code, 200)
}
