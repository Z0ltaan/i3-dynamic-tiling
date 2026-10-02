package main

import (
	"log"
	i3 "go.i3wm.org/i3/v4"
)

func setLayout(event *i3.WindowEvent) {
	command := ""
	if event.Container.Rect.Width <= event.Container.Rect.Height {
		command = "split v"
	} else {
		command = "split h"
	}
	_, err := i3.RunCommand(command)
	if err != nil {
		log.Printf("Error while running command: %v", err)
	}
}

func isFocusEvent(event *i3.WindowEvent) bool {
	return event.Change == "focus" || event.Change == "new"
}

func main() {

	listener := i3.Subscribe("window")
	defer listener.Close()

	log.Println("listening")
	for listener.Next() {
		log.Println("New event")
		ev, ok := listener.Event().(*i3.WindowEvent)
		if !ok {
			continue
		}

		if isFocusEvent(ev) {
			log.Println("Focused Window")
			setLayout(ev)
		}
	}

}
