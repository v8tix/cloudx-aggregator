package store

import (
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/model/request"
	"github.com/samber/lo"
	"sync"
)

type (
	AssociationsStore struct {
		Associations                  []request.Association
		LastAssociationsCorrelationID string
		LastMessageCorrelationID      string
		mu                            sync.Mutex
	}

	ParentGroup struct {
		Source      string
		Destination string
	}
)

func NewAssociationsStore() *AssociationsStore {
	return &AssociationsStore{
		Associations: make([]request.Association, 0),
	}
}

func newParentGroup(source, destination string) ParentGroup {
	return ParentGroup{
		Source:      source,
		Destination: destination,
	}
}

func (a *AssociationsStore) Add(group *dto.GroupDTO) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.LastAssociationsCorrelationID != group.AssociationsDTO.Metadata.CorrelationID {
		a.Associations = append(a.Associations, *group.AssociationsDTO.Associations...)
		a.LastAssociationsCorrelationID = group.AssociationsDTO.Metadata.CorrelationID
	}
}

func (a *AssociationsStore) FindParentsByChildren(group *dto.GroupDTO) (ParentGroup, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var source, destination string

	if a.LastMessageCorrelationID == group.MessageDTO.Metadata.CorrelationID {
		return ParentGroup{}, false
	}

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

	isOk := len(source) != 0 && len(destination) != 0

	if isOk {
		a.LastMessageCorrelationID = group.MessageDTO.Metadata.CorrelationID
	}

	return newParentGroup(source, destination), isOk
}
