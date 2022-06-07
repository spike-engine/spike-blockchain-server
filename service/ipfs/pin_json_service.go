package ipfs

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/model"
	"spike-blockchain-server/serializer"
)

type PinJsonService struct {
	Json string `json:"json"`
}

func (service *PinJsonService) PinJson() serializer.Response {
	config := model.DefaultPinataConfig
	config.PinataContent = service.Json
	options, err := json.Marshal(config.PinataOptions)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}

	metadata, err := json.Marshal(config.PinataMetaData)
	if err != nil {
		return serializer.Response{
			Code:  401,
			Error: err.Error(),
		}
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("pinata_api_key", os.Getenv("PINATA_API_KEY")).
		SetHeader("pinata_secret_api_key", os.Getenv("PINATA_SECRET_KEY")).
		SetFormData(map[string]string{
			"pinataContent":  config.PinataContent,
			"pinataOptions":  string(options),
			"pinataMetadata": string(metadata),
		}).
		Post(PINATA_PIN_JSON)
	if err != nil {
		return serializer.Response{
			Code:  402,
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return serializer.Response{
			Code:  403,
			Error: resp.String(),
		}
	}

	var res serializer.Pin
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return serializer.Response{
			Code:  405,
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: 200,
		Data: res,
	}
}
