package ipfs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/go-resty/resty/v2"

	"spike-blockchain-server/model"
	"spike-blockchain-server/serializer"
)

type PinJsonFileService struct {
	File *multipart.FileHeader `form:"file" json:"file" binding:"required"`
	Name string                `form:"name" json:"name"`
}

func (service *PinJsonFileService) PinJsonFile() serializer.Response {
	config := model.DefaultPinataConfig
	config.PinataMetadata.Name = service.File.Filename
	options, err := json.Marshal(config.PinataOptions)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}

	metadata, err := json.Marshal(config.PinataMetadata)
	if err != nil {
		return serializer.Response{
			Code:  401,
			Error: err.Error(),
		}
	}

	file, err := service.File.Open()
	if err != nil {
		return serializer.Response{
			Code:  402,
			Error: err.Error(),
		}
	}

	all, err := ioutil.ReadAll(file)
	if err != nil {
		return serializer.Response{
			Code:  403,
			Error: err.Error(),
		}
	}
	var out bytes.Buffer
	err = json.Indent(&out, all, "", "\t")
	if err != nil {
		return serializer.Response{
			Code:  404,
			Error: err.Error(),
		}
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("pinata_api_key", os.Getenv("PINATA_API_KEY")).
		SetHeader("pinata_secret_api_key", os.Getenv("PINATA_SECRET_KEY")).
		SetFileReader("file", service.File.Filename, bytes.NewReader(out.Bytes())).
		SetFormData(map[string]string{
			"pinataOptions":  string(options),
			"pinataMetadata": string(metadata),
		}).
		Post(PINATA_PIN_FILE)
	if err != nil {
		return serializer.Response{
			Code:  405,
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return serializer.Response{
			Code:  406,
			Error: resp.String(),
		}
	}

	var res serializer.Pin
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return serializer.Response{
			Code:  407,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: res,
	}
}
