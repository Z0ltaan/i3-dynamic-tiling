package server

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Z0ltaan/i3-tiling-hook/internal/split"
	"go.i3wm.org/i3/v4"
)

type LayoutServer struct {
	stop         atomic.Bool
	config       ServerConfig
	listener     net.Listener
	listenerLock sync.Mutex
}

func NewLayoutServer(cfg ServerConfig) *LayoutServer {
	return &LayoutServer{config: cfg}
}

func (server *LayoutServer) SetUpNextSplit(event *i3.WindowEvent) {
	server.listenerLock.Lock()
	server.config.strategy.SetUpNextSplit(event)
	server.listenerLock.Unlock()
}

func (server *LayoutServer) Close() {
	server.stop.Store(true)
	os.Remove(server.config.socketPath)
}

func (s *LayoutServer) Listen() error {
	var err error
	s.listener, err = net.Listen("unix", s.config.socketPath)
	if err != nil {
		return err
	}
	defer s.listener.Close()

	for !s.stop.Load() {
		log.Println("Waiting for new connections on socket")
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error while accepting connection: ", err)
			continue
		}

		buffer, err := io.ReadAll(conn)
		if err != nil {
			log.Println("Error while reading from socket: ", err)
			continue
		}
		commandAndData := strings.Split(string(buffer), ":")
		if len(commandAndData) != 2 {
			log.Println("Wrong command format")
			continue
		}

		command, data := commandAndData[0], commandAndData[1]
		if err := s.runCommand(command, data); err != nil {
			log.Println("Failed to run a command: ", err)
			continue
		}

	}

	return nil
}

type UnknownCommandError struct {
	cause string
}

func (e *UnknownCommandError) Error() string {
	return e.cause
}

func (server *LayoutServer) runCommand(cmd string, data string) error {
	commands := map[string]func(string) error{
		"strategy": func(s string) error { return server.handleSplitStrategyChange(s) },
	}

	dataHandler, ok := commands[cmd]
	if !ok {
		return &UnknownCommandError{cause: cmd + " is an unknown command"}
	}

	return dataHandler(data)
}

func (server *LayoutServer) handleSplitStrategyChange(newStratName string) error {
	strat, err := getSplitStrategy(newStratName)
	if err != nil {
		return err
	}

	server.config.strategy = strat
	return nil
}

func getSplitStrategy(stratName string) (split.SplitStrategy, error) {
	strategies := map[string]func() split.SplitStrategy{
		"smart":      func() split.SplitStrategy { return split.NewSmartSplitStrategy() },
		"horizontal": func() split.SplitStrategy { return split.NewHorizontalSplitStrategy() },
		"vertical":   func() split.SplitStrategy { return split.NewVerticalSplitStrategy() },
		// "spiral":     func() split.SplitStrategy {return split.NewSpiralStrategy() },
		"off": func() split.SplitStrategy { return split.NewOffSplitStrategy() },
	}

	contructor, ok := strategies[stratName]
	if !ok {
		return nil, &UnknownCommandError{cause: "Unknown split strategy"}
	}

	return contructor(), nil
}
