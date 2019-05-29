package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

// TestFFProbeStreamFormat ...
func TestFFProbeStreamFormat(t *testing.T) {
	format, _ := FFProbeStreamFormat("D:\\video\\周杰伦唱歌贼难听.mp4")
	v, _ := json.Marshal(format)
	ioutil.WriteFile("d:\\test.log", v, os.ModePerm)
	t.Log(string(v))
}
