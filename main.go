package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Size(1280, 720))
		
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

	editingArea := newEditingArea()

	colorPicker := newColorPicker(f32.Pt(0, 0), image.Pt(1280, 50))

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			
			title := material.H1(theme, fmt.Sprintf("%v", colorPicker.chosenColor))
			title.Color = colorPicker.chosenColor
			title.Alignment = text.Middle
			title.Layout(gtx)

			editingArea.board.drawingColor = colorPicker.chosenColor
			
			editingArea.Layout(gtx)

			colorPicker.size.X = e.Size.X
			colorPicker.position.Y = float32(e.Size.Y - colorPicker.size.Y)
			colorPicker.Layout(gtx)

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
