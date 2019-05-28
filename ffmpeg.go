package cmd

import (
	"context"
)

// FFMpegSpliteMedia ...
func FFMpegSpliteMedia(ctx context.Context, args string) {
	ffmpeg := NewFFMpeg()
	ffmpeg.SetArgs(args)
	info := make(chan string, 1024)
	cls := make(chan bool)

	go func() {
		e := ffmpeg.RunContext(ctx, info, cls)
		if e != nil {
			return
		}
	}()
	for {
		select {
		case v := <-info:
			if v != "" {
				log.With("status", "process").Info(v)
			}
		case c := <-cls:
			if c == true {
				close(info)
				return
			}
		case <-ctx.Done():
			log.With("status", "done")
			if err := ctx.Err(); err != nil {
				log.Error(err)
			}
			return
		default:
			//log.Println("waiting:...")
		}
	}
}
