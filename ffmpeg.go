package cmd

import (
	"context"
	"fmt"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const sliceM3u8FFmpegTemplate = "-y -i %s -strict -2  -c:v %s -c:a %s -bsf:v h264_mp4toannexb -f hls -hls_list_size 0 -hls_time %d  -hls_segment_filename %s %s"

// SplitArgs ...
type SplitArgs struct {
	Auto            bool
	Output          string
	Video           string
	Audio           string
	M3U8            string
	SegmentFileName string
	HLSTime         int
}

// FFmpegContext ...
type ffmpegContext struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// Context ...
func (c *ffmpegContext) Context() context.Context {
	return c.ctx
}

// Add ...
func (c *ffmpegContext) Add(i int) {
	c.wg.Add(i)
}

// Wait ...
func (c *ffmpegContext) Wait() {
	c.wg.Wait()
}

// Done ...
func (c *ffmpegContext) Done() {
	c.wg.Done()
}

// Cancel ...
func (c *ffmpegContext) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Context ...
type Context interface {
	Cancel()
	Add(int)
	Wait()
	Done()
	Context() context.Context
}

// FFmpegContext ...
func FFmpegContext() Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &ffmpegContext{
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
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

// AutoOption ...
func AutoOption(s bool) SplitOptions {
	return func(args *SplitArgs) {
		args.Auto = s
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

// FFMpegSplitToM3U8WithProbe ...
func FFMpegSplitToM3U8WithProbe(ctx Context, file string, args ...SplitOptions) (e error) {
	if strings.Index(file, " ") != -1 {
		return xerrors.New("file name cannot have spaces")
	}
	sa := SplitArgs{
		Output:          "",
		Auto:            true,
		Video:           "libx264",
		Audio:           "aac",
		M3U8:            "media.m3u8",
		SegmentFileName: "media-%05d.ts",
		HLSTime:         10,
	}
	for _, o := range args {
		o(&sa)
	}
	format, e := FFProbeStreamFormat(file)
	if e != nil {
		return e
	}
	video := format.Video()
	audio := format.Audio()
	if !format.IsVideo() || audio == nil || video == nil {
		return xerrors.New("open file failed with ffprobe")
	}
	if video.CodecName == sa.Video {
		sa.Video = "copy"
	}
	if audio.CodecName == sa.Audio {
		sa.Audio = "copy"
	}
	out, e := filepath.Abs(sa.Output)
	if e != nil {
		return e
	}

	if sa.Auto {
		_ = os.MkdirAll(out, os.ModePerm)
	}

	sfn := filepath.Join(out, sa.SegmentFileName)
	m3u8 := filepath.Join(out, sa.M3U8)

	tpl := fmt.Sprintf(sliceM3u8FFmpegTemplate, file, sa.Video, sa.Audio, sa.HLSTime, sfn, m3u8)
	return FFMpegRun(ctx, tpl)
}

// FFMpegSplitToM3U8 ...
func FFMpegSplitToM3U8(ctx Context, file string, args ...SplitOptions) (e error) {
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

	out, e := filepath.Abs(sa.Output)
	if e != nil {
		return e
	}
	//_ = os.MkdirAll(out, os.ModePerm)

	sfn := filepath.Join(out, sa.SegmentFileName)
	m3u8 := filepath.Join(out, sa.M3U8)

	tpl := fmt.Sprintf(sliceM3u8FFmpegTemplate, file, sa.Video, sa.Audio, sa.HLSTime, sfn, m3u8)
	return FFMpegRun(ctx, tpl)
}

// FFMpegRun ...
func FFMpegRun(ctx Context, args string) (e error) {
	ffmpeg := NewFFMpeg()
	ffmpeg.SetArgs(args)
	info := make(chan string, 1024)
	go func() {
		e = ffmpeg.RunContext(ctx, info)
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
		case <-ctx.Context().Done():
			log.With("status", "done")
			if e = ctx.Context().Err(); e != nil {
				if e == context.Canceled {
					log.Info("exit with cancel")
				}
			}
			return
		default:
			//log.Println("waiting:...")
		}
	}
}
