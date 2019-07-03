package cmd

import (
	"testing"
)

// TestFormat_NameAnalyze ...
func TestFormat_NameAnalyze(t *testing.T) {
	sf1, e := FFProbeStreamFormat("d:\\video\\周杰伦唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4")

	if e != nil {
		log.Error(e)
		return
	}
	t.Logf("%+v", sf1.Video())
}
