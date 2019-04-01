package cmd

import (
	"context"
	"log"
	"testing"
)

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	ctx := context.Background()
	s := make(chan string, 1024)
	b := make(chan bool)
	command := Default()
	cmd := command.Input("D:\\video\\12.28.mp4").
		Split("D:\\video\\output")

	go func() {
		e := cmd.RunContext(ctx, s, b)
		//(ctx, s, b)
		if e != nil {
			log.Println("error:", e)
		}
	}()
	for {
		select {
		case v := <-s:
			if v != "" {
				log.Print(v)
			}
		case c := <-b:
			if c == true {
				close(s)
				return
			}
		case <-ctx.Done():
			log.Println("done")
			log.Println(ctx.Err())
			return
		default:
			//log.Println("waiting:...")
		}
	}
}
