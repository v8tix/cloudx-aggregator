package pipe

import (
	"encoding/json"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/model/request"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
)

type AssociationsObservable struct {
	logger *slog.Logger
}

func NewAssociationsObservable(logger *slog.Logger) AssociationsObservable {
	return AssociationsObservable{logger: logger}
}

func (a AssociationsObservable) Pipe(obs rxgo.Observable) rxgo.Observable {
	filter := filterByType(obs, a.areAssociations)
	associationsMap := mapBytesTo(filter, toMany[request.Association])
	associationDTOMap := mapToMany[request.Association, dto.AssociationsDTO](associationsMap, a.toAssociationsDTO)
	return associationDTOMap
}

func (a AssociationsObservable) areAssociations(data []uint8) bool {
	var associations []request.Association
	if err := json.Unmarshal(data, &associations); err == nil {
		return true
	}
	return false
}

func (a AssociationsObservable) toAssociations(data []uint8) ([]request.Association, error) {
	var associations []request.Association

	err := json.Unmarshal(data, &associations)
	if err != nil {
		return nil, err
	}
	return associations, nil
}

func (a AssociationsObservable) toAssociationsDTO(associations *[]request.Association) *dto.AssociationsDTO {
	return dto.NewAssociationsDTO(associations)
}
