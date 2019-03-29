package cmd

type Command struct {
	Name string
	Opts map[string][]string
}

func New(name string) *Command {
	return &Command{
		Name: name,
		Opts: make(map[string][]string),
	}
}

func (c *Command) Ignore() *Command {
	c.Opts["ignore"] = []string{"-y"}
	return c
}

func (c *Command) Input(path string) *Command {
	c.Opts["input"] = []string{"-i", path}
	return c
}

func (c *Command) Codec() *Command {
	c.Opts["c"] = []string{"-c", "copy"}
	return c
}

type String func(out *string)

func Strings(in string) String {
	return func(out *string) {
		*out = in
	}
}

func (c *Command) CodecVideo(options ...String) *Command {
	option := "copy"
	for _, v := range options {
		v(&option)
	}
	c.Opts["c"] = []string{"-c:v", option}
	return c
}
func (c *Command) CodecAudio(options ...String) *Command {
	option := "copy"
	for _, v := range options {
		v(&option)
	}
	c.Opts["c"] = []string{"-c:a", option}
	return c
}
