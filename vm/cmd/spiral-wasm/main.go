package main

import (
	"image"
	"log"
	"syscall/js"

	"libdb.so/vm/internal/hypnospiral"
)

var imgBuf = make([]byte, 1000*1000) // 1 MB buffer

func main() {
	global := js.Global()
	global.Set("hypnospiral_draw", js.FuncOf(draw))
	global.Set("hypnospiral_ready", js.ValueOf(true))

	select {}
}

// draw(buffer: Uint8Array, width: number, height: number, t: number, params: { spinSpeed: number, zoom: number, blur: number })
func draw(this js.Value, args []js.Value) any {
	b := args[0]       // Uint8Array buffer
	w := args[1].Int() // width
	h := args[2].Int() // height
	t := float64(args[3].Float()) / 1000
	jsParams := args[4]

	pixLen := w * h
	if pixLen > len(imgBuf) {
		log.Printf("hypnospiral: refusing to draw, image size %dx%d exceeds buffer capacity %d", w, h, len(imgBuf))
		return js.ValueOf(false)
	}

	img := &image.Gray{
		Pix:    imgBuf[:pixLen:pixLen],
		Rect:   image.Rect(0, 0, w, h),
		Stride: w,
	}

	params := hypnospiral.Parameters{
		SpinSpeed: jsParams.Get("spinSpeed").Float(),
		Zoom:      jsParams.Get("zoom").Float(),
		Blur:      jsParams.Get("blur").Float(),
	}

	hypnospiral.Draw(img, params, t)

	ok := js.CopyBytesToJS(b, img.Pix) == pixLen
	return js.ValueOf(ok)
}
