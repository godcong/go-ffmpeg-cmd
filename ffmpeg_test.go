package cmd

import "testing"

// TestFFMpegRun ...
func TestFFMpegRun(t *testing.T) {

	ctx := FFmpegContext()
	//ctx.Add(1)

	go FFMpegSplitToM3U8WithProbe(ctx, "/mnt/d/video/女大学生的沙龙室.Room.Salon.College.Girls.2018.HD720P.X264.AAC.Korean.CHS.mp4", OutputOption("tmp"))
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
