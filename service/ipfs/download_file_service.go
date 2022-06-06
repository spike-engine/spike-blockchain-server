package ipfs

import (
	"github.com/cavaliergopher/grab/v3"
	"spike-blockchain-server/serializer"
)

type DownloadFileService struct {
	SaveFilePath string `json:"saveFilePath"`
	CID          string `json:"cid"`
}

func (service *DownloadFileService) DownloadFile() serializer.Response {
	client := grab.NewClient()
	req, err := grab.NewRequest(service.SaveFilePath, DOWNLOAD_ADDR+service.CID)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}
	resp := client.Do(req)
	if err = resp.Err(); err != nil {
		return serializer.Response{
			Code:  resp.HTTPResponse.StatusCode,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: *service,
	}
}
