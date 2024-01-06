package rx

import (
	"context"
	"encoding/json"
	"fmt"
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
	return obs.Filter(func(item interface{}) bool {

		data, ok := item.([]uint8)
		if !ok {
			a.logger.Error(fmt.Errorf("unexpected type: %T, expected []uint8", item).Error())
			return false
		}

		return a.areAssociations(data)

	}).Map(func(_ context.Context, item interface{}) (interface{}, error) {

		data, ok := item.([]uint8)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected []uint8", item)
			a.logger.Error(err.Error())
			return nil, err
		}

		return a.toAssociations(data)

	}).Map(func(_ context.Context, item interface{}) (interface{}, error) {

		data, ok := item.([]request.Association)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected []request.Association", item)
			a.logger.Error(err.Error())
			return nil, err
		}

		return a.toAssociationsDTO(data), nil
	})
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

func (a AssociationsObservable) toAssociationsDTO(associations []request.Association) *dto.AssociationsDTO {
	return dto.NewAssociationsDTO(associations)
}
