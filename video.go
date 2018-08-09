package main

import (
	"image"
	"log"
	"time"

	"github.com/3d0c/gmf"
)

func (app *tdApp) startVideoCB(s string, i interface{}) {

	var err error

	app.videoChan, err = drone.VideoConnectDefault()
	if err != nil {
		alertDialog(app, errorSev, err.Error())
	}

	// start video feed when drone connects
	drone.StartVideo()
	go func() {
		for {
			drone.StartVideo()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	app.videoStopChan = make(chan bool) // unbuffered

	go app.videoListener()
}

func (app *tdApp) customReader() ([]byte, int) {
	block := <-app.videoChan
	return block, len(block)
}

func assert(i interface{}, err error) interface{} {
	if err != nil {
		log.Fatalf("Assert %v", err)
	}

	return i
}

func (app *tdApp) videoListener() {

	iCtx := gmf.NewCtx()
	defer iCtx.CloseInputAndRelease()

	if err := iCtx.SetInputFormat("h264"); err != nil {
		log.Fatalf("iCtx SetInputFormat %v", err)
	}

	avioCtx, err := gmf.NewAVIOContext(iCtx, &gmf.AVIOHandlers{ReadPacket: app.customReader})
	defer gmf.Release(avioCtx)
	if err != nil {
		log.Fatalf("NewAVIOContext %v", err)
	}

	iCtx.SetPb(avioCtx).OpenInput("")
	if err != nil {
		log.Fatalf("iCtx OpenInput %v", err)
	}

	srcVideoStream, err := iCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Fatalf("GetBestStream %v", err)
	}

	// codec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_PNG)
	codec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_RAWVIDEO)
	if err != nil {
		log.Fatalf("FindDecoder %v", err)
	}
	cc := gmf.NewCodecCtx(codec)
	defer gmf.Release(cc)

	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	// cc.SetPixFmt(gmf.AV_PIX_FMT_RGB24).
	cc.SetPixFmt(gmf.AV_PIX_FMT_BGR32).
		SetWidth(1280).
		SetHeight(720).
		SetTimeBase(gmf.AVR{Num: 1, Den: 1})

	if err := cc.Open(nil); err != nil {
		log.Fatalf("cc Open %v", err)
	}

	swsCtx := gmf.NewSwsCtx(srcVideoStream.CodecCtx(), cc, gmf.SWS_BICUBIC)
	defer gmf.Release(swsCtx)

	dstFrame := gmf.NewFrame().
		SetWidth(1280).
		SetHeight(720).
		SetFormat(gmf.AV_PIX_FMT_BGR32) //SetFormat(gmf.AV_PIX_FMT_RGB32)
	defer gmf.Release(dstFrame)

	if err := dstFrame.ImgAlloc(); err != nil {
		log.Fatalf("ImgAlloc %v", err)
	}

	ist := assert(iCtx.GetStream(srcVideoStream.Index())).(*gmf.Stream)
	defer gmf.Release(ist)

	codecCtx := ist.CodecCtx()
	defer gmf.Release(codecCtx)

	for pkt := range iCtx.GetNewPackets() {

		if pkt.StreamIndex() != srcVideoStream.Index() {
			app.Log().Info("Skipping wrong stream packet")
			continue
		}

		frame, err := pkt.Frames(codecCtx)
		if err != nil {
			app.Log().Info("CodeCtx %v", err)
			continue
		}

		swsCtx.Scale(frame, dstFrame)

		p, err := dstFrame.Encode(cc)
		if err != nil {
			app.Log().Fatal("Encode %v", err)
		}

		rgba := new(image.RGBA)
		rgba.Stride = 4 * 1280
		rgba.Rect = image.Rect(0, 0, 1280, 720)
		rgba.Pix = p.Data()

		app.texture.SetFromRGBA(rgba)
		app.feed.SetChanged(true)
		app.Log().Info("Frame decoded")

		gmf.Release(p)
		gmf.Release(frame)
		gmf.Release(pkt)

	}
}
