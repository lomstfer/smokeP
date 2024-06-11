package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"smokep/utils"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/sqweek/dialog"
)

var g_theme *material.Theme

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Size(1280, 720), app.Title("smokeP"))
		window.Option(app.NavigationColor(color.NRGBA{0, 0, 0, 255}))

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
	g_theme.Bg = color.NRGBA{50, 50, 50, 255}
	g_theme.Fg = color.NRGBA{255, 255, 255, 255}
	var ops op.Ops

	editingArea := newEditingArea()
	settingsArea := newSettingsArea()
	go func() {
		for {
			select {
			case <-settingsArea.SaveButtonClicked:
				filePath, err := dialog.File().Title("").Filter("PNG image file", "png").SetStartFile("export.png").Save()
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					utils.SaveImageToFile(editingArea.board.pixelImg, filePath)
				}
			case <-settingsArea.LoadButtonClicked:
				filePath, err := dialog.File().Title("").Filter("PNG image file", "png").Load()
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					img := utils.LoadImage(filePath)
					if img == nil {
						fmt.Println("rip")
						break
					}
					editingArea.board.setToNewImage(img)
					window.Invalidate()
				}
			}
		}
	}()

	background := paint.NewImageOp(utils.GenerateGridImage(160, 90, color.NRGBA{200, 200, 200, 255}, color.NRGBA{100, 100, 100, 255}))
	background.Filter = paint.FilterNearest

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			{
				r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
				area := clip.Rect(r).Push(gtx.Ops)

				background.Add(gtx.Ops)

				scale := max(float32(r.Dx())/float32(background.Size().X), float32(r.Dy())/float32(background.Size().Y))
				tStack := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				tStack.Pop()

				area.Pop()
			}

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return settingsArea.Layout(g_theme, gtx)
				}),
				layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
					editingArea.board.drawingColor = settingsArea.colorPicker.ChosenColor
					return editingArea.Layout(gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
