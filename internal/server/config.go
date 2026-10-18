package server

import (
	"github.com/Z0ltaan/i3-tiling-hook/internal/split"
	// "os"
)

type ServerConfig struct {
	strategy   split.SplitStrategy
	socketPath string
	// cmd handler
}

type SocketError struct {
	cause string
}

func (e *SocketError) Error() string {
	return e.cause
}

func NewTilerConfig(setUps ...func(*ServerConfig) error) (ServerConfig, error) {
	config := ServerConfig{}
	for _, setUp := range setUps {
		if err := setUp(&config); err != nil {
			return ServerConfig{}, err
		}
	}

	return config, nil
}

func WithSocketPath(path string) func(*ServerConfig) error {
	// os.Stat()
	return func(cfg *ServerConfig) error { cfg.socketPath = path; return nil }
}

func WithStrategy(stratName string) func(*ServerConfig) error {
	return func(cfg *ServerConfig) error {
		strat, err := getSplitStrategy(stratName)
		if err != nil {
			return err
		}
		cfg.strategy = strat
		return nil
	}
}
