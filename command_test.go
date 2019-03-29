package cmd

import "testing"

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	command := Default()
	s, e := command.Input("D:\\video\\12.28.mp4").
		Output("D:\\video\\output\\media.m3u8").
		Run()
	t.Log(s, e)
}
