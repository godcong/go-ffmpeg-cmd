package cmd

import (
	"bufio"
	"github.com/godcong/go-trait"
	"io"
	"io/ioutil"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var log = trait.NewZapSugar()

// Command ...
type Command struct {
	Path string
	Name string
	Args []string
	//OutPath string
	//Opts    map[string][]string
}

// New ...
func New(name string) *Command {
	return &Command{
		Name: name,
		//Opts: make(map[string][]string),
	}
}

// NewFFMpeg ...
func NewFFMpeg() *Command {
	return New("ffmpeg")
}

// NewFFProbe ...
func NewFFProbe() *Command {
	return New("ffprobe")
}

// SetArgs ...
func (c *Command) SetArgs(s string) {
	c.Args = strings.Split(s, " ")
}

// AddArgs ...
func (c *Command) AddArgs(s string) {
	c.Args = append(c.Args, s)
}

// SetPath ...
func (c *Command) SetPath(s string) {
	c.Path = s
}

// CMD ...
func (c *Command) CMD() string {
	if c.Path != "" {
		return filepath.Join(c.Path, c.Name)
	}
	return c.Name
}

// GetCurrentDir ...
func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return dir
}

// Run ...
func (c *Command) Run() (string, error) {
	cmd := exec.Command(c.CMD(), c.Args...)
	cmd.Env = c.Env()
	//显示运行的命令
	log.With("run", "Run").Info(cmd.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}
	return string(stdout), nil
}

// Env ...
func (c *Command) Env() []string {
	path := os.Getenv("PATH")
	if err := os.Setenv("PATH", path+":"+GetCurrentDir()); err != nil {
		//err = xerrors.Errorf("PATH error:%+v", err)
		log.Error(err)
	}
	return os.Environ()
}

// RunContext ...
func (c *Command) RunContext(ctx Context, info chan<- string) (e error) {
	cmd := exec.CommandContext(ctx.Context(), c.CMD(), c.Args...)

	//显示运行的命令
	log.With("run", "RunContext").Info(cmd.Args)
	defer func() {
		ctx.Done()
		if e != nil {
			log.Error(e)
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

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		select {
		case <-ctx.Context().Done():
			return ctx.Context().Err()
		default:
			lines, _, e := reader.ReadLine()
			log.Info(string(lines), e)
			if e != nil || io.EOF == e {
				goto END
			}
			if strings.TrimSpace(string(lines)) != "" {
				if info != nil {
					info <- string(lines)
				}
			}
		}
	}
END:
	if e != nil {
		bytes, _ := ioutil.ReadAll(stderr)
		info <- string(bytes)
	}
	e = cmd.Wait()
	if e != nil {
		return e
	}

	return nil
}
