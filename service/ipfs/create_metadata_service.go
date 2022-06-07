package ipfs

import (
	"encoding/json"
	"spike-blockchain-server/model"
	"spike-blockchain-server/serializer"
)

type CreateMetaDataService struct {
	MetaJson []byte
}

func (service *CreateMetaDataService) CreateMetaData() serializer.Response {
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
