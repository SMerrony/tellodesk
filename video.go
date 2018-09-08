package main

import (
	"bufio"
	"image"
	"log"
	"os"
	"time"

	"github.com/mattn/go-gtk/glib"

	"github.com/3d0c/gmf"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

func buildFeedWgt() (wgt *gtk.Image) {
	wgt = gtk.NewImageFromFile(bluesky)
	return wgt
}

func recordVideoCB() {
	var vidPath string
	fs := gtk.NewFileChooserDialog(
		"Save Video Recording to...",
		win,
		gtk.FILE_CHOOSER_ACTION_SAVE, "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
	fs.SetCurrentFolder(settings.DataDir)
	ff := gtk.NewFileFilter()
	ff.AddPattern("*.h264")
	fs.SetFilter(ff)
	res := fs.Run()
	if res == gtk.RESPONSE_ACCEPT {
		vidPath = fs.GetFilename()
		if vidPath != "" {
			var err error
			videoFile, err = os.OpenFile(vidPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				alert := gtk.NewMessageDialog(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_CLOSE,
					"Could not create video file.")
				alert.SetTitle(appName)
				alert.Run()
				alert.Destroy()
			} else {
				videoWriter = bufio.NewWriter(videoFile)
				videoRecMu.Lock()
				videoRecording = true
				videoRecMu.Unlock()
				menuBar.recVidItem.SetSensitive(false)
				menuBar.stopRecVidItem.SetSensitive(true)
			}
		}
	}
	fs.Destroy()
}

func stopRecordingVideoCB() {
	videoRecMu.Lock()
	videoRecording = false
	videoRecMu.Unlock()
	videoWriter.Flush()
	videoFile.Close()
	menuBar.recVidItem.SetSensitive(true)
	menuBar.stopRecVidItem.SetSensitive(false)
}

func startVideo() {

	var err error

	videoChan, err = drone.VideoConnectDefault()
	if err != nil {
		log.Print(err.Error())
		alert := gtk.NewMessageDialog(win, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, err.Error())
		alert.SetTitle(appName)
		alert.Run()
		alert.Destroy()
	}

	// start video feed restarter when drone connects
	drone.StartVideo()
	go func() { // no GTK stuff in here...
		for {
			drone.StartVideo()
			select {
			case <-vrStopChan:
				return
			default:
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	stopNewPicChan = make(chan bool)
	newPicChan = make(chan bool)
	go videoListener()

	//go updateFeed()

	glib.TimeoutAdd(30, func() bool {
		updateFeed()
		return true
	})
}

func customReader() ([]byte, int) {
	pkt := <-videoChan
	videoRecMu.RLock()
	if videoRecording {
		videoWriter.Write(pkt)
	}
	videoRecMu.RUnlock()
	return pkt, len(pkt)
}

func assert(i interface{}, err error) interface{} {
	if err != nil {
		log.Fatalf("Assert %v", err)
	}

	return i
}

//func (app *tdApp) videoListener() {
func videoListener() {

	//Log().Info("Videolistener started")

	iCtx := gmf.NewCtx()
	defer iCtx.CloseInputAndRelease()

	if err := iCtx.SetInputFormat("h264"); err != nil {
		log.Fatalf("iCtx SetInputFormat %v", err)
	}
	//Log().Info("Input format set")
	avioCtx, err := gmf.NewAVIOContext(iCtx, &gmf.AVIOHandlers{ReadPacket: customReader})
	defer gmf.Release(avioCtx)
	if err != nil {
		log.Fatalf("NewAVIOContext %v", err)
	}

	//Log().Info("Setting Pb...")
	iCtx.SetPb(avioCtx)

	//Log().Info("Opening input...")
	err = iCtx.OpenInput("")
	if err != nil {
		log.Fatalf("iCtx OpenInput %v", err)
	}

	//Log().Info("Getting best stream...")
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
		SetWidth(videoWidth).
		SetHeight(videoHeight).
		SetTimeBase(gmf.AVR{Num: 1, Den: 1})
	//Log().Info("Opening cc")
	if err := cc.Open(nil); err != nil {
		log.Fatalf("cc Open %v", err)
	}

	swsCtx := gmf.NewSwsCtx(srcVideoStream.CodecCtx(), cc, gmf.SWS_BICUBIC)
	defer gmf.Release(swsCtx)

	dstFrame := gmf.NewFrame().
		SetWidth(videoWidth).
		SetHeight(videoHeight).
		SetFormat(gmf.AV_PIX_FMT_BGR32) //SetFormat(gmf.AV_PIX_FMT_RGB32)
	defer gmf.Release(dstFrame)

	if err := dstFrame.ImgAlloc(); err != nil {
		log.Fatalf("ImgAlloc %v", err)
	}

	ist := assert(iCtx.GetStream(srcVideoStream.Index())).(*gmf.Stream)
	defer gmf.Release(ist)

	codecCtx := ist.CodecCtx()
	defer gmf.Release(codecCtx)

	//Log().Info("Entering get video packets loop...")

	for pkt := range iCtx.GetNewPackets() {

		if pkt.StreamIndex() != srcVideoStream.Index() {
			log.Println("Skipping wrong stream packet")
			continue
		}

		frame, err := pkt.Frames(codecCtx)
		if err != nil {
			log.Printf("CodeCtx %v", err)
			continue
		}

		swsCtx.Scale(frame, dstFrame)

		p, err := dstFrame.Encode(cc)

		if err != nil {
			log.Fatalf("Encode %v", err)
		}
		rgba := new(image.RGBA)
		rgba.Stride = 4 * videoWidth
		rgba.Rect = image.Rect(0, 0, videoWidth, videoHeight)
		rgba.Pix = p.Data()

		newFeedImageMu.Lock()
		feedImage = rgba
		newFeedImage = true
		newFeedImageMu.Unlock()

		gmf.Release(p)
		gmf.Release(frame)
		gmf.Release(pkt)

	}
}

// updateFeed actually updates the video image in the feed tab.
// It must be run on the main thread, so there is a little mutex dance to
// check if a new image is ready for display.
func updateFeed() {
	newFeedImageMu.Lock()
	if newFeedImage {
		var pbd gdkpixbuf.PixbufData
		pbd.Colorspace = gdkpixbuf.GDK_COLORSPACE_RGB
		pbd.HasAlpha = true
		pbd.BitsPerSample = 8
		pbd.Width = videoWidth
		pbd.Height = videoHeight
		pbd.RowStride = videoWidth * 4 // RGBA

		pbd.Data = feedImage.Pix

		pb := gdkpixbuf.NewPixbufFromData(pbd)
		feedWgt.SetFromPixbuf(pb)

		newFeedImage = false
	}
	newFeedImageMu.Unlock()
}
