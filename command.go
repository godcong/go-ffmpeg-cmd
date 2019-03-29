package cmd

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

// String ...
type String func(out *string)

// Strings ...
func Strings(in string) String {
	return func(out *string) {
		*out = in
	}
}

// CodecVideo ...
func (c *Command) CodecVideo(options ...String) *Command {
	option := "copy"
	for _, v := range options {
		v(&option)
	}
	c.Opts["cv"] = []string{"-c:v", option}
	return c
}

// CodecAudio ...
func (c *Command) CodecAudio(options ...String) *Command {
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
	c.Opts["hls_segment_filename"] = []string{"-hls_segment_filename", n}
	return c
}
