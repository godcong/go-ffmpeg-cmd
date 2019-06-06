package cmd

import (
	"os"
	"testing"
)

// TestFormat_NameAnalyze ...
func TestFormat_NameAnalyze(t *testing.T) {
	sf1, e := FFProbeStreamFormat("d:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4")

	if e != nil {
		log.Error(e)
		return
	}
	t.Logf("%+v", sf1.Format)
	_ = os.Rename("d:\\video\\周杰伦唱歌贼难听.mp4", sf1.NameAnalyze().ToString())
	analyze := sf1.NameAnalyze().ToString()
	t.Log(analyze)

	sf2, e := FFProbeStreamFormat("d:\\video\\[BT天堂btbttt.com]我的女友.My.Girlfriend.2018.HD720P.X264.AAC.Korean.中文字幕.mp4")
	if e != nil {
		return
	}
	t.Log(sf2.Format)
	analyze2 := sf2.NameAnalyze()
	t.Logf("%+v", *analyze2)

}
