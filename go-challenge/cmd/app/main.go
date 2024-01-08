package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/configuration"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/pipe"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ch := make(chan rxgo.Item)

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
	groupObs := pipe.NewGroupObservable(logger)

	go pipe.Producer(logger, conn, ch)

	rawObservable := rxgo.FromChannel(ch)
	associationsPipe := associationsObs.Pipe(rawObservable)
	messagePipe := messageObs.Pipe(rawObservable)
	groupPipe := groupObs.Pipe(associationsPipe, messagePipe)

	for item := range groupPipe.Observe() {
		group := item.V.(*dto.GroupDTO)
		bytes, _ := json.Marshal(group)
		fmt.Printf("value: %s\n", string(bytes))
		if pipe.AssociationAggregator.FindParentByChildren(group) {
			logger.Info("Parent found!")
		}
		fmt.Printf("value:%#v\n", item.V)
	}

	select {}
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
	logger.Info("attempting to connect", "server:", uri)
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	logger.Info("connection established", "server:", uri)
	return conn, nil
}

// TODO give a name for the output binary
func printUsage() {
	fmt.Println("Usage: cloudx-client [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
