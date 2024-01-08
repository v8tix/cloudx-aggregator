package main

import (
	"flag"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/configuration"
	"github.com/cloudx-labs/challenge/internal/pipe"
	"github.com/cloudx-labs/challenge/internal/store"
	"github.com/cloudx-labs/challenge/internal/task"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var done chan struct{}

func main() {
	var wg sync.WaitGroup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("client", "process", os.Getpid())
	ch := make(chan rxgo.Item)
	done = make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	wg.Add(1)
	go func() {
		<-interrupt
		wg.Done()
		close(done)
	}()
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	associationStore := store.NewAssociationsStore()
	responseStore := store.NewResponseStore()

	conn, err := run(logger)
	if err != nil {
		logger.Error(err.Error())
	}

	defer func(conn *websocket.Conn) {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}(conn)

	associationsObs := pipe.NewAssociationsObservable(logger)
	messageObs := pipe.NewMessageObservable(logger)
	groupObs := pipe.NewGroupObservable(logger, associationStore)
	rawObservable := rxgo.FromChannel(ch)
	associationsPipe := associationsObs.Pipe(rawObservable)
	messagePipe := messageObs.Pipe(rawObservable)
	groupPipe := groupObs.Pipe(associationsPipe, messagePipe)
	groupDTOCh := groupPipe.Observe()

	wg.Add(1)
	go task.Producer(logger, conn, ch, done, &wg)
	wg.Add(1)
	go task.Write(logger, responseStore, done, &wg)
	wg.Add(1)
	go task.Aggregator(logger, associationStore, responseStore, groupDTOCh, done, &wg)
	wg.Wait()
}

func run(logger *slog.Logger) (*websocket.Conn, error) {
	var cfg configuration.RemoteServerCfg

	flag.StringVar(&cfg.Host, "host", "localhost", "specify the remote server's address (default is localhost)")
	flag.StringVar(&cfg.Port, "port", "5050", "remote port to listen on for requests")
	showHelp := flag.Bool("help", false, "show help message")
	flag.Parse()

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	uri := fmt.Sprintf("ws://%s:%s", cfg.Host, cfg.Port)
	logger.Info("attempting to connect", "server", uri)
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	logger.Info("connection established", "server", uri)
	return conn, nil
}

// TODO give a name for the output binary
func printUsage() {
	fmt.Println("Usage: cloudx-client [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
