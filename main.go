package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()

	var ops op.Ops

	b := newEditingArea()

	app.Decorated(false)

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			title := material.H1(theme, "smoke")
			title.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
			title.Alignment = text.Middle
			title.Layout(gtx)

			grect := image.Rect(100, 100, 1000, 200)
			paint.LinearGradientOp{
				Stop1:  f32.Pt(float32(grect.Min.X), float32(grect.Min.Y)),
				Stop2:  f32.Pt(float32(grect.Max.X - grect.Dx()/2), float32(grect.Max.Y)),
				Color1: color.NRGBA{255, 0, 0, 255},
				Color2: color.NRGBA{255, 255, 0, 255},
			}.Add(gtx.Ops)
			paint.LinearGradientOp{
				Stop1:  f32.Pt(float32(grect.Min.X + grect.Dx()/2), float32(grect.Min.Y)),
				Stop2:  f32.Pt(float32(grect.Max.X), float32(grect.Max.Y)),
				Color1: color.NRGBA{255, 255, 0, 255},
				Color2: color.NRGBA{0, 0, 255, 255},
			}.Add(gtx.Ops)
			garea := clip.Rect(grect).Push(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			garea.Pop()

			b.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

// func getPixelData(img image.Image) *image.NRGBA {
// 	rgba := image.NewNRGBA(img.Bounds())

// 	for y := 0; y < img.Bounds().Dx(); y++ {
// 		for x := 0; x < img.Bounds().Dy(); x++ {
// 			c := img.At(x, y)
// 			rgba.Set(x, y, c)
// 		}
// 	}

// 	return rgba
// }
