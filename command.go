package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Command ...
type Command struct {
	Name    string
	OutPath string
	Opts    map[string][]string
}

// New ...
func New(name string, path string) *Command {
	return &Command{
		Name:    name,
		OutPath: path,
		Opts:    make(map[string][]string),
	}
}

// Default ...
func Default() *Command {
	return New("ffmpeg", "").
		Ignore().CodecAudio().CodecVideo().
		BitStreamFiltersVideo("h264_mp4toannexb").Format("hls").HlsTime("10").
		HlsListSize("0")
}

// Split ...
func (c *Command) Split(path string) *Command {
	return c.SetPath(path).HlsSegmentFilename("media-%03d.ts").Output("media.m3u8")
}

// SetPath ...
func (c *Command) SetPath(path string) *Command {
	c.OutPath = path
	return c
}

// Ignore ...
func (c *Command) Ignore() *Command {
	c.Opts["ignore"] = []string{"-y"}
	return c
}

// Input ...
func (c *Command) Input(path string) *Command {
	c.Opts["input"] = []string{"-i", path}
	return c
}

// Codec ...
func (c *Command) Codec() *Command {
	c.Opts["c"] = []string{"-c", "copy"}
	return c
}

// Strings ...
type Strings func(out *string)

// String ...
func String(in string) Strings {
	return func(out *string) {
		*out = in
	}
}

// CodecVideo ...
func (c *Command) CodecVideo(options ...Strings) *Command {
	option := "copy"
	for _, v := range options {
		v(&option)
	}
	c.Opts["cv"] = []string{"-c:v", option}
	return c
}

// CodecAudio ...
func (c *Command) CodecAudio(options ...Strings) *Command {
	option := "copy"
	for _, v := range options {
		v(&option)
	}
	c.Opts["ca"] = []string{"-c:a", option}
	return c
}

// Format ...
func (c *Command) Format(f string) *Command {
	c.Opts["format"] = []string{"-f", f}
	return c
}

// HlsTime ...
func (c *Command) HlsTime(t string) *Command {
	c.Opts["hls_time"] = []string{"-hls_time", t}
	return c
}

// HlsListSize ...
func (c *Command) HlsListSize(s string) *Command {
	c.Opts["hls_list_size"] = []string{"-hls_list_size", s}
	return c
}

// HlsSegmentFilename ...
func (c *Command) HlsSegmentFilename(name string) *Command {
	_, file := filepath.Split(name)
	c.Opts["hls_segment_filename"] = []string{"-hls_segment_filename", filepath.Join(c.OutPath, file)}
	return c
}

// HlsKeyInfoFile ...
func (c *Command) HlsKeyInfoFile(file string) *Command {
	c.Opts["hls_key_info_file"] = []string{"-hls_key_info_file", file}
	return c
}

// BitStreamFiltersVideo ...
func (c *Command) BitStreamFiltersVideo(f string) *Command {
	c.Opts["bsfv"] = []string{"-bsf:v", f}
	return c
}

// Output ...
func (c *Command) Output(name string) *Command {
	_, file := filepath.Split(name)
	c.Opts["output"] = []string{filepath.Join(c.OutPath + file)}
	return c
}

// Options ...
func (c *Command) Options() []string {
	var options []string
	input, b := c.Opts["input"]
	if !b {
		return nil
	}
	output, b := c.Opts["output"]
	if !b {
		return nil
	}
	delete(c.Opts, "input")
	delete(c.Opts, "output")
	options = append(options, input...)
	for _, v := range c.Opts {
		options = append(options, v...)
	}
	options = append(options, output...)
	return options
}

// Run ...
func (c *Command) Run() (string, error) {
	cmd := exec.Command(c.Name, c.Options()...)
	cmd.Env = os.Environ()
	fmt.Println(cmd.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}
	return string(stdout), nil
}

// RunContext ...
func (c *Command) RunContext(ctx context.Context, s chan<- string, close chan<- bool) (e error) {
	cmd := exec.CommandContext(ctx, c.Name, c.Options()...)
	cmd.Env = os.Environ()
	//显示运行的命令
	log.Println("command:", cmd.Args)
	defer func() {
		close <- true
		if e != nil {
			panic(e)
		}
	}()
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		return e
	}

	stderr, e := cmd.StderrPipe()
	if e != nil {
		return e
	}

	e = cmd.Start()
	if e != nil {
		return e
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))
	//实时循环读取输出流中的一行内容
	for {
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			log.Println("end", cmd.Args, e)
			break
		}
		if line != "" {
			s <- "lineinfo:" + line
		}
	}

	e = cmd.Wait()
	if e != nil {
		return e
	}
	return nil
}
