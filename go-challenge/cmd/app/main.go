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
	logger.Info("client", "pid", os.Getpid())
	readerCh := make(chan rxgo.Item)
	done = make(chan struct{})
	interrupt := make(chan os.Signal, 1)

	conn, err := connect(logger)
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

	associationStore := store.NewAssociationsStore()
	responseStore := store.NewResponseStore()
	associationsObs := pipe.NewAssociationsObservable(logger)
	messageObs := pipe.NewMessageObservable(logger)
	groupObs := pipe.NewGroupObservable(logger, associationStore)
	rawObservable := rxgo.FromChannel(readerCh)
	associationsPipe := associationsObs.Pipe(rawObservable)
	messagePipe := messageObs.Pipe(rawObservable)
	groupPipe := groupObs.Pipe(associationsPipe, messagePipe)
	groupDTOCh := groupPipe.Observe()

	tasks := task.NewTasks(
		logger,
		responseStore,
		associationStore,
		conn,
		readerCh,
		groupDTOCh,
		interrupt,
		done,
		&wg,
	)

	tasks.Run()
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	wg.Wait()
}

func connect(logger *slog.Logger) (*websocket.Conn, error) {
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

func printUsage() {
	fmt.Println("Usage: wsclient [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
