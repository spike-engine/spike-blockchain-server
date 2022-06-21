package ipfs

import (
	"os"

	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/serializer"
)

type DownloadService struct {
	IpfsHash string `form:"ipfs_hash" json:"ipfs_hash" binding:"required"`
}

func (service *DownloadService) Download() serializer.Response {
	client := resty.New()

	resp, err := client.R().Get(service.url())
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}

	if resp.IsError() {
		return serializer.Response{
			Code:  401,
			Error: resp.String(),
		}
	}
	return serializer.Response{
		Code: 200,
		Data: resp.Body(),
		Msg:  resp.Header().Get("Content-Type"),
	}
}

func (service *DownloadService) url() string {
	return PROTOCOL + "://" + os.Getenv("PINATA_GATEWAY") + "/ipfs/" + service.IpfsHash
}
