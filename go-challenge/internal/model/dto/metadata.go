package dto

import (
	"github.com/google/uuid"
	"time"
)

type Metadata struct {
	CorrelationId string
	ReceivedAt    time.Time
}

func NewMetadata() Metadata {
	return Metadata{
		CorrelationId: uuid.NewString(),
		ReceivedAt:    time.Now().UTC(),
	}
}
