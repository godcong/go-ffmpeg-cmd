package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -ss %s -to %s -c:v %s -c:a %s -bsf:v h264_mp4toannexb -vsync 0 -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -c:v %s -c:a %s -bsf:v h264_mp4toannexb -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`

// SplitArgs ...
type SplitArgs struct {
	StreamFormat    *StreamFormat
	Auto            bool
	Start           string
	End             string
	Output          string
	Video           string
	Audio           string
	M3U8            string
	SegmentFileName string
	HLSTime         int
	probe           func(string) (*StreamFormat, error)
}

// FFmpegContext ...
type ffmpegContext struct {
	once   sync.Once
	mu     sync.RWMutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	done   chan bool
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
	select {
	case <-c.Waiting():
		return
	}
}

// Waiting ...
func (c *ffmpegContext) Waiting() <-chan bool {
	c.once.Do(func() {
		go func() {
			c.wg.Wait()
			c.done <- true
		}()
	})

	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan bool)
	}
	d := c.done
	c.mu.Unlock()
	return d
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
	Waiting() <-chan bool
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

// StreamFormatOption ...
func StreamFormatOption(s *StreamFormat) SplitOptions {
	return func(args *SplitArgs) {
		args.StreamFormat = s
	}
}

// ProbeInfoOption ...
func ProbeInfoOption(f func(string) (*StreamFormat, error)) SplitOptions {
	return func(args *SplitArgs) {
		args.probe = f
	}
}

// FFMpegSplitToM3U8WithProbe ...
func FFMpegSplitToM3U8WithProbe(ctx Context, file string, args ...SplitOptions) (sa *SplitArgs, e error) {
	args = append(args, ProbeInfoOption(FFProbeStreamFormat))
	return FFMpegSplitToM3U8(ctx, file, args...)
}

// FFMpegSplitToM3U8 ...
func FFMpegSplitToM3U8(ctx Context, file string, args ...SplitOptions) (sa *SplitArgs, e error) {
	if strings.Index(file, " ") != -1 {
		return nil, xerrors.New("file name cannot have spaces")
	}
	if ctx == nil {
		ctx = FFmpegContext()
	}
	sa = &SplitArgs{
		Output:          "",
		Auto:            true,
		Video:           "libx264",
		Audio:           "aac",
		M3U8:            "media.m3u8",
		SegmentFileName: "media-%05d.ts",
		HLSTime:         10,
	}
	for _, o := range args {
		o(sa)
	}

	if sa.probe != nil {
		sa.StreamFormat, e = sa.probe(file)
		if e != nil {
			return nil, e
		}
	}
	if sa.StreamFormat != nil {
		video := sa.StreamFormat.Video()
		audio := sa.StreamFormat.Audio()
		if !sa.StreamFormat.IsVideo() || audio == nil || video == nil {
			return nil, xerrors.New("open file failed with ffprobe")
		}
		if video.CodecName == "h264" {
			sa.Video = "copy"
		}
		if audio.CodecName == sa.Audio {
			sa.Audio = "copy"
		}

		if sa.Start == "" {
			sa.Start = video.StartTime
		}
		if sa.End == "" {
			sa.End = video.Duration
		}

	}

	sa.Output, e = filepath.Abs(sa.Output)
	if e != nil {
		return nil, e
	}
	log.With("output", sa.Output).Info("output dir")
	if sa.Auto {
		sa.Output = filepath.Join(sa.Output, uuid.New().String())
		_ = os.MkdirAll(sa.Output, os.ModePerm)
	}

	sfn := filepath.Join(sa.Output, sa.SegmentFileName)
	m3u8 := filepath.Join(sa.Output, sa.M3U8)

	tpl := fmt.Sprintf(sliceM3u8FFmpegTemplate, file, sa.Video, sa.Audio, sa.HLSTime, sfn, m3u8)
	time.Sleep(time.Duration(rand.Int31n(30)) * time.Second)
	return nil, FFMpegRun(ctx, tpl)
}

// FFMpegRun ...
func FFMpegRun(ctx Context, args string) (e error) {
	ffmpeg := NewFFMpeg()
	ffmpeg.SetArgs(args)
	info := make(chan string, 1024)
	done := make(chan error, 1)
	go func() {
		ctx.Add(1)
		done <- ffmpeg.RunContext(ctx, info)
	}()
	for {
		select {
		case e = <-done:
			if e != nil {
				log.Error(e)
			}
			return
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
