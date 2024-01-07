package dto

import (
	"github.com/google/uuid"
	"time"
)

type Metadata struct {
	CorrelationID string
	ReceivedAt    time.Time
}

func NewMetadata() Metadata {
	return Metadata{
		CorrelationID: uuid.NewString(),
		ReceivedAt:    time.Now().UTC(),
	}
}
