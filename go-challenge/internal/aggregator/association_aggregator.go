package aggregator

import (
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/model/request"
	"github.com/samber/lo"
	"sync"
)

type AssociationsAggregator struct {
	Associations      []request.Association
	LastCorrelationID string
	mu                sync.Mutex
}

func NewAssociationsAggregator() *AssociationsAggregator {
	return &AssociationsAggregator{
		Associations:      make([]request.Association, 0),
		LastCorrelationID: "",
	}
}

func (a *AssociationsAggregator) AddAssociations(group *dto.GroupDTO) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.LastCorrelationID != group.AssociationsDTO.Metadata.CorrelationID {
		a.Associations = append(a.Associations, *group.AssociationsDTO.Associations...)
		a.LastCorrelationID = group.AssociationsDTO.Metadata.CorrelationID
	}
}

func (a *AssociationsAggregator) FindParentByChildren(group *dto.GroupDTO) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	var source, destination string

	aMap := lo.Associate(a.Associations, func(f request.Association) (string, string) {
		return f.Children, f.Parent
	})

	sourceParent, ok := aMap[group.MessageDTO.Message.Source]
	if ok {
		source = sourceParent
	}

	destinationParent, ok := aMap[group.MessageDTO.Message.Destination]
	if ok {
		destination = destinationParent
	}

	return len(source) != 0 && len(destination) != 0
}
