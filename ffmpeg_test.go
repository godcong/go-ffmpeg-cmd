package cmd

import "testing"

// TestFFMpegRun ...
func TestFFMpegRun(t *testing.T) {

	ctx := FFmpegContext()
	ctx.Add(1)

	e := FFMpegSplitToM3U8(ctx, "D:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4")
	t.Error(e)
}
