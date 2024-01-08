package task

import (
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/store"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Tasks struct {
	logger           *slog.Logger
	store            *store.ResponseStore
	associationStore *store.AssociationsStore
	conn             *websocket.Conn
	readerCh         chan<- rxgo.Item
	groupDTOCh       <-chan rxgo.Item
	interrupt        <-chan os.Signal
	done             chan struct{}
	wg               *sync.WaitGroup
	wsMutex          sync.Mutex
}

func NewTasks(
	logger *slog.Logger,
	store *store.ResponseStore,
	associationStore *store.AssociationsStore,
	conn *websocket.Conn,
	next chan<- rxgo.Item,
	groupDTOCh <-chan rxgo.Item,
	interrupt <-chan os.Signal,
	done chan struct{},
	wg *sync.WaitGroup,
) Tasks {
	return Tasks{
		logger:           logger,
		store:            store,
		associationStore: associationStore,
		conn:             conn,
		readerCh:         next,
		groupDTOCh:       groupDTOCh,
		interrupt:        interrupt,
		done:             done,
		wg:               wg,
	}
}

func (t *Tasks) Run() {
	t.wg.Add(4)
	go t.signal()
	go t.wsReader()
	go t.wsWriter()
	go t.aggregator()
}

func (t *Tasks) signal() {
	<-t.interrupt
	t.wg.Done()
	t.logger.Info("Signal", "msg:", "shutdown")
	close(t.done)
}

func (t *Tasks) wsWriter() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bytes, err := t.store.ToResponse()
			if err != nil {
				t.logger.Error(err.Error())
				continue
			}

			t.logger.Info("writing", "store: ", string(bytes))
			t.wsMutex.Lock()
			err = t.conn.WriteMessage(websocket.TextMessage, bytes)
			t.wsMutex.Unlock()
			if err != nil {
				t.logger.Error(err.Error())
				return
			}

		case <-t.done:
			t.wg.Done()
			t.logger.Info("Write", "msg:", "shutdown")
			return
		}
	}
}

func (t *Tasks) wsReader() {
	for {
		select {
		case <-t.done:
			t.wg.Done()
			t.logger.Info("Producer", "msg:", "shutdown")
			return
		default:
			t.wsMutex.Lock()
			_, message, err := t.conn.ReadMessage()
			t.wsMutex.Unlock()
			if err != nil {
				t.logger.Error(err.Error())
				return
			}
			t.readerCh <- rxgo.Of(message)
		}
	}
}

func (t *Tasks) aggregator() {
	for {
		select {
		case <-t.groupDTOCh:

			value := <-t.groupDTOCh
			group := value.V.(*dto.GroupDTO)
			parentGroup, ok := t.associationStore.FindParentsByChildren(group)
			if ok {
				t.store.Add(parentGroup)
			}

		case <-t.done:
			t.wg.Done()
			t.logger.Info("Aggregator", "msg:", "shutdown")
			return
		}
	}
}
