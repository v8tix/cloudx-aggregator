package pipe

import (
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/store"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
)

type GroupObservable struct {
	AssociationStore *store.AssociationsStore
	logger           *slog.Logger
}

func NewGroupObservable(
	logger *slog.Logger,
	associationStore *store.AssociationsStore,
) GroupObservable {
	return GroupObservable{
		AssociationStore: associationStore,
		logger:           logger,
	}
}

func (g *GroupObservable) Pipe(observables ...rxgo.Observable) rxgo.Observable {
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
			g.AssociationStore.Add(group)
			return group
		},
		observables,
	)
}
