package model

import (
	"errors"
)

type SpikeMetadata struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	ExternalURL string       `json:"external_url"`
	SpikeInfo   SpikeInfo    `json:"spike_info"`
	Attributes  []Attributes `json:"attributes"`
}

type SpikeInfo struct {
	Version        string `json:"version"`
	SpikeModelURL  string `json:"spike_model_url"`
	SpikeModelType string `json:"spike_model_type"`
}

type Attributes struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

func (service *SpikeMetadata) ValidateMetaData() error {
	if len(service.Name) == 0 || len(service.Name) >= 10 {
		return errors.New("metadata name is Illegal")
	}

	if len(service.Description) == 0 {
		return errors.New("Description cannot be empty ")
	}

	if len(service.SpikeInfo.SpikeModelType) == 0 || len(service.SpikeInfo.Version) == 0 || len(service.SpikeInfo.SpikeModelURL) == 0 {
		return errors.New("SpikeInfo cannot be empty ")
	}

	for _, attribute := range service.Attributes {
		if len(attribute.Value) == 0 || len(attribute.TraitType) == 0 {
			return errors.New("attributes is Illegal")
		}
	}
	return nil
}
