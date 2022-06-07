package ipfs

import (
	"encoding/json"

	"spike-blockchain-server/model"
	"spike-blockchain-server/serializer"
)

type CreateMetadataService struct {
	MetaJson []byte `json:"meta_json"`
}

func (service *CreateMetadataService) CreateMetadata() serializer.Response {
	metaData := model.SpikeMetaData{}
	err := json.Unmarshal(service.MetaJson, &metaData)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Error: err.Error(),
		}
	}
	err = metaData.ValidateMetaData()
	if err != nil {
		return serializer.Response{
			Code:  401,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: string(service.MetaJson),
	}
}
