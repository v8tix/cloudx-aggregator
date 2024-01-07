package dto

import (
	"github.com/cloudx-labs/challenge/internal/model/request"
)

type AssociationsDTO struct {
	Metadata
	Associations *[]request.Association
}

func (a AssociationsDTO) isDTO() {}

func NewAssociationsDTO(associations *[]request.Association) *AssociationsDTO {
	return &AssociationsDTO{
		Metadata:     NewMetadata(),
		Associations: associations,
	}
}
