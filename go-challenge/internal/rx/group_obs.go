package rx

import (
	"github.com/cloudx-labs/challenge/internal/aggregator"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
)

var (
	AssociationAggregator = aggregator.NewAssociationsAggregator()
)

type GroupObservable struct {
	logger *slog.Logger
}

func NewGroupObservable(logger *slog.Logger) GroupObservable {
	return GroupObservable{logger: logger}
}

func (g GroupObservable) Pipe(observables ...rxgo.Observable) rxgo.Observable {
	return rxgo.CombineLatest(
		func(i ...interface{}) interface{} {
			var associationsDTO *dto.AssociationsDTO
			var messageDTO *dto.MessageDTO
			var err error

			for _, v := range i {
				switch v := v.(type) {
				case *dto.AssociationsDTO:
					associationsDTO = v
				case *dto.MessageDTO:
					messageDTO = v
				default:
					err = v.(error)
				}
			}
			if err != nil {
				return err
			}

			group := dto.NewGroup(associationsDTO, messageDTO)
			AssociationAggregator.AddAssociations(group)
			return group
		},
		observables,
	)
}
