package cmd

import (
	"context"
	"log"
	"testing"
)

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	s := make(chan string)
	b := make(chan bool)
	command := Default()
	cmd := command.Input("D:\\video\\12.28.mp4").
		Split("D:\\video\\output")
	go cmd.RunContext(context.Background(), s, b)

	for {
		select {
		case v := <-s:
			log.Println("loading", v)
		case c := <-b:
			if c == true {
				close(s)
				break
			}
		}
	}

}
