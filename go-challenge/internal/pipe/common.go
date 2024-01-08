package pipe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/model/request"
	"github.com/reactivex/rxgo/v2"
)

func filterByType(obs rxgo.Observable, f func([]uint8) bool) rxgo.Observable {
	return obs.Filter(func(item interface{}) bool {
		data, ok := item.([]uint8)
		if !ok {
			return false
		}

		return f(data)
	})
}

func mapBytesTo(obs rxgo.Observable, f func(data []uint8) (interface{}, error)) rxgo.Observable {
	return obs.Map(func(_ context.Context, item interface{}) (interface{}, error) {
		data, ok := item.([]uint8)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected []uint8", item)
			return nil, err
		}

		return f(data)
	})
}

func mapTo[R request.ReqI, T dto.DTOI](obs rxgo.Observable, f func(data *R) T) rxgo.Observable {
	return obs.Map(func(_ context.Context, item interface{}) (interface{}, error) {
		data, ok := item.(*R)
		if !ok {
			err := fmt.Errorf("unexpected type: %T", item)
			return nil, err
		}

		return f(data), nil
	})
}

func mapToMany[R request.ReqI, T dto.DTOI](obs rxgo.Observable, f func(data *[]R) *T) rxgo.Observable {
	return obs.Map(func(_ context.Context, item interface{}) (interface{}, error) {
		data, ok := item.(*[]R)
		if !ok {
			err := fmt.Errorf("unexpected type: %T", item)
			return nil, err
		}

		return f(data), nil
	})
}

func to[T request.ReqI](data []uint8) (interface{}, error) {
	var req T

	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func toMany[T request.ReqI](data []uint8) (interface{}, error) {
	var req []T

	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
