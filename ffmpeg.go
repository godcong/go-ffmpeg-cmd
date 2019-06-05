package cmd

import (
	"context"
	"fmt"
	"golang.org/x/xerrors"
	"strings"
)

const sliceM3u8FFmpegTemplate = "-y -i %s -strict -2  -c:v %s -c:a %s -bsf:v h264_mp4toannexb -f hls -hls_list_size 0 -hls_time %d  -hls_segment_filename %s %s"

// SplitArgs ...
type SplitArgs struct {
	Output          string
	Video           string
	Audio           string
	M3U8            string
	SegmentFileName string
	HLSTime         int
}

// SplitOptions ...
type SplitOptions func(args *SplitArgs)

// HLSTimeOption ...
func HLSTimeOption(i int) SplitOptions {
	return func(args *SplitArgs) {
		args.HLSTime = i
	}
}

// OutputOption ...
func OutputOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Output = s
	}
}

// VideoOption ...
func VideoOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Video = s
	}
}

// AudioOption ...
func AudioOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Video = s
	}
}

// FFMpegSplitToM3U8 ...
func FFMpegSplitToM3U8(file string, args ...SplitOptions) (err error) {
	if strings.Index(file, " ") != -1 {
		return xerrors.New("file name cannot have spaces")
	}
	sa := SplitArgs{
		Output:          "",
		Video:           "libx264",
		Audio:           "aac",
		M3U8:            "media.m3u8",
		SegmentFileName: "media-%05d.ts",
		HLSTime:         10,
	}
	for _, o := range args {
		o(&sa)
	}
	tpl := fmt.Sprintf(sliceM3u8FFmpegTemplate, file, sa.Video, sa.Audio, sa.HLSTime, sa.SegmentFileName, sa.M3U8)
	FFMpegRun(context.Background(), tpl)
	return nil
}

// FFMpegRun ...
func FFMpegRun(ctx context.Context, args string) {
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
