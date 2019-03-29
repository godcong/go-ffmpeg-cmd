package cmd

import "testing"

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	command := New("ffmpeg")
	s, e := command.Input("test.mp4").
		Ignore().CodecAudio(String("aac")).CodecVideo(String("libx264")).
		BitStreamFiltersVideo("h264_mp4toannexb").Format("hls").HlsTime("10").
		HlsListSize("0").HlsSegmentFilename("media-%03d.ts").
		Output("media.m3u8").
		Run()
	t.Log(s, e)
}
