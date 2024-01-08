package task

import (
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/store"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
	"sync"
	"time"
)

func Write(logger *slog.Logger, store *store.ResponseStore, done chan struct{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bytes, err := store.ToResponse()
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			logger.Info("write", "store: ", string(bytes))

		case <-done:
			wg.Done()
			logger.Info("Write", "msg:", "shutdown")
			return
		}
	}
}

func Producer(
	logger *slog.Logger,
	conn *websocket.Conn,
	next chan<- rxgo.Item,
	done <-chan struct{},
	wg *sync.WaitGroup,
) {
	for {
		select {
		case <-done:
			wg.Done()
			logger.Info("Producer", "msg:", "shutdown")
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Error(err.Error())
				return
			}
			next <- rxgo.Of(message)
		}
	}
}

func Aggregator(
	logger *slog.Logger,
	associationStore *store.AssociationsStore,
	store *store.ResponseStore,
	groupDTOCh <-chan rxgo.Item,
	done chan struct{},
	wg *sync.WaitGroup,
) {
	for {
		select {
		case <-groupDTOCh:

			value := <-groupDTOCh
			group := value.V.(*dto.GroupDTO)
			parentGroup, ok := associationStore.FindParentsByChildren(group)
			if ok {
				store.Add(parentGroup)
			}

		case <-done:
			wg.Done()
			logger.Info("Aggregator", "msg:", "shutdown")
			return
		}
	}
}
