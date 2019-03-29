package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Command ...
type Command struct {
	Name string
	Opts map[string][]string
}

// New ...
func New(name string) *Command {
	return &Command{
		Name: name,
		Opts: make(map[string][]string),
	}
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
	c.Opts["hls_segment_filename"] = []string{"-hls_segment_filename", name}
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
func (c *Command) Output(path string) *Command {
	c.Opts["output"] = []string{path}
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
