package pipe

import (
	"encoding/json"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/model/request"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
)

type MessageObservable struct {
	logger *slog.Logger
}

func NewMessageObservable(logger *slog.Logger) MessageObservable {
	return MessageObservable{logger: logger}
}

func (m MessageObservable) Pipe(obs rxgo.Observable) rxgo.Observable {
	filter := filterByType(obs, m.isMessage)
	messageMap := mapBytesTo(filter, to[request.Message])
	messageDTOMap := mapTo[request.Message](messageMap, m.toMessageDTO)
	return messageDTOMap
}

func (m MessageObservable) isMessage(data []uint8) bool {
	var message request.Message

	if err := json.Unmarshal(data, &message); err == nil {
		return true
	}
	return false
}

func (m MessageObservable) toMessageDTO(message *request.Message) *dto.MessageDTO {
	return dto.NewMessageDTO(message)
}
