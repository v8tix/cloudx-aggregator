package dto

import (
	"github.com/cloudx-labs/challenge/internal/model/request"
)

type AssociationsDTO struct {
	Metadata
	associations []request.Association
}

func NewAssociationsDTO(associations []request.Association) *AssociationsDTO {
	return &AssociationsDTO{
		Metadata:     NewMetadata(),
		associations: associations,
	}
}
