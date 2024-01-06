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

type MessageObservable struct {
	logger *slog.Logger
}

func NewMessageObservable(logger *slog.Logger) MessageObservable {
	return MessageObservable{logger: logger}
}

func (m MessageObservable) Pipe(obs rxgo.Observable) rxgo.Observable {
	return obs.Filter(func(item interface{}) bool {

		data, ok := item.([]uint8)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected []uint8", item)
			m.logger.Error(err.Error())
			return false
		}

		return m.isMessage(data)

	}).Map(func(_ context.Context, item interface{}) (interface{}, error) {

		data, ok := item.([]uint8)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected []uint8", item)
			m.logger.Error(err.Error())
			return nil, err
		}

		return m.toMessage(data)

	}).Map(func(_ context.Context, item interface{}) (interface{}, error) {

		data, ok := item.(*request.Message)
		if !ok {
			err := fmt.Errorf("unexpected type: %T, expected request.Message", item)
			m.logger.Error(err.Error())
			return nil, err
		}

		return m.toMessageDto(data), nil
	})
}

func (m MessageObservable) isMessage(data []uint8) bool {
	var message request.Message

	if err := json.Unmarshal(data, &message); err == nil {
		return true
	}
	return false
}

func (m MessageObservable) toMessage(data []uint8) (*request.Message, error) {
	var message request.Message

	err := json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (m MessageObservable) toMessageDto(message *request.Message) dto.MessageDTO {
	return dto.NewMessageDTO(message)
}
