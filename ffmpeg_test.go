package cmd

import "testing"

// TestFFMpegRun ...
func TestFFMpegRun(t *testing.T) {

	ctx := FFmpegContext()
	ctx.Add(1)

	go FFMpegSplitToM3U8WithProbe(ctx, "D:\\video\\周杰伦唱歌贼难听.2019.1080P.h264.aac.Japanese.None..mp4", OutputOption("tmp"))
	//ctx.Add(1)
	//go FFMpegSplitToM3U8(ctx, "D:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4", OutputOption("tmp"))
	//ctx.Add(1)
	//go FFMpegSplitToM3U8(ctx, "D:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4", OutputOption("tmp"))
	//ctx.Add(1)
	//go FFMpegSplitToM3U8(ctx, "D:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4", OutputOption("tmp"))
	//ctx.Add(1)
	//go FFMpegSplitToM3U8(ctx, "D:\\video\\极乐女忍者.LADY.NINJA.2018.HD1080P.X264.AAC.Japanese.CHT.mp4", OutputOption("tmp"))
	//for {
	ctx.Wait()
	//}

}
