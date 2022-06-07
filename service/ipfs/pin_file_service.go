package ipfs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/model"
	"spike-blockchain-server/serializer"
)

type PinFileService struct {
	Filepath string `json:"filepath"`
}

func (service *PinFileService) PinFile() serializer.Response {
	config := model.DefaultPinataConfig
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

	f, err := ioutil.ReadFile(service.Filepath)
	if err != nil {
		return serializer.Response{
			Code:  402,
			Error: err.Error(),
		}
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("pinata_api_key", os.Getenv("PINATA_API_KEY")).
		SetHeader("pinata_secret_api_key", os.Getenv("PINATA_SECRET_KEY")).
		SetFileReader("file", ".env.example", bytes.NewReader(f)).
		SetFormData(map[string]string{
			"pinataOptions":  string(options),
			"pinataMetadata": string(metadata),
		}).
		Post(PINATA_PIN_FILE)
	if err != nil {
		return serializer.Response{
			Code:  403,
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return serializer.Response{
			Code:  404,
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
