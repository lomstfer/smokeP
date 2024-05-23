package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

var g_theme *material.Theme

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Size(640, 360))

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	g_theme = material.NewTheme()

	var ops op.Ops

	editingArea := newEditingArea()

	settingsArea := newSettingsArea()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return settingsArea.Layout(gtx)
				}),
				layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
					editingArea.board.drawingColor = settingsArea.colorPicker.chosenColor
					return editingArea.Layout(gtx)
				}),
			)

			grect := image.Rect(0, 0, 200, 200)
			{
				paint.LinearGradientOp{
					Stop1:  f32.Pt(float32(grect.Min.X), float32(grect.Min.Y)),
					Stop2:  f32.Pt(float32(grect.Max.X), float32(grect.Min.Y)),
					Color1: color.NRGBA{255, 255, 255, 255},
					Color2: color.NRGBA{255, 0, 0, 255},
				}.Add(gtx.Ops)
				garea := clip.Rect(grect).Push(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				garea.Pop()
			}

			{
				paint.LinearGradientOp{
					Stop1:  f32.Pt(float32(grect.Max.X), float32(grect.Max.Y)),
					Stop2:  f32.Pt(float32(grect.Max.X), float32(grect.Min.Y)),
					Color1: color.NRGBA{0, 0, 0, 255},
					Color2: color.NRGBA{0, 0, 0, 0},
				}.Add(gtx.Ops)
				garea := clip.Rect(grect).Push(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				garea.Pop()
			}

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
