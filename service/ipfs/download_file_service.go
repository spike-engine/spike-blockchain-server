package ipfs

import (
	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/serializer"
)

type DownloadFileService struct {
}

func (service *DownloadFileService) DownloadFile() serializer.Response {
	client := resty.New()
	_ = client

	return serializer.Response{
		Code: 200,
		Data: nil,
	}
}
