package dto

import "github.com/cloudx-labs/challenge/internal/model/request"

type MessageDTO struct {
	Metadata
	Message *request.Message
}

func NewMessageDTO(message *request.Message) *MessageDTO {
	return &MessageDTO{
		Metadata: NewMetadata(),
		Message:  message,
	}
}
