package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/Z0ltaan/i3-tiling-hook/internal/server"
	i3 "go.i3wm.org/i3/v4"
)

func isFocusEvent(event *i3.WindowEvent) bool {
	return event.Change == "focus" || event.Change == "new"
}

func main() {
	config, err := server.NewTilerConfig(
		server.WithSocketPath("/tmp/i3dt.socket"),
		server.WithStrategy("smart"),
	)

	if err != nil {
		log.Println("Error constructing a config: ", err)
		return
	}

	layoutServer := server.NewLayoutServer(config)
	go layoutServer.Listen()

	tilingStop := &atomic.Bool{}
	signalHandler := func(server *server.LayoutServer) {
		channel := make(chan os.Signal, 2)
		signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
		for sig := range channel {
			server.Close()
			tilingStop.Store(true)
			log.Println("Recieved ", sig.String(), ". Closing...")
		}
	}

	go signalHandler(layoutServer)
	listener := i3.Subscribe("window")
	defer listener.Close()

	log.Println("listening")
	for listener.Next() && !tilingStop.Load() {
		ev, ok := listener.Event().(*i3.WindowEvent)
		if !ok {
			continue
		}

		if isFocusEvent(ev) {
			log.Println("New focus event")
			layoutServer.SetUpNextSplit(ev)
		}
	}

}
