package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"smokep/utils"

	"gioui.org/app"
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
	settingsArea := newSettingsArea(editingArea.board.pixelImgOp.Size())
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
						fmt.Println("Error loading image")
						break
					}
					editingArea.board.pixelImg = img
					editingArea.board.refreshImage()
					editingArea.board.centerImage()
					window.Invalidate()
				}
			case newSize := <-settingsArea.PixelBoardSizeEditor.editorSubmit:
				editingArea.board.Resize(newSize, settingsArea.PixelBoardSizeEditor.selectedOrigin)
			}
		}
	}()

	gridBackground := utils.NewGridBackground()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			gridBackground.WindowSize = gtx.Constraints.Max
			editingArea.Update(gtx)
			settingsArea.Update(gtx, editingArea.board.pixelImgOp.Size())

			{ // background color
				r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
				area := clip.Rect(r).Push(gtx.Ops)
				paint.ColorOp{Color: g_theme.Bg}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				area.Pop()
			}

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return settingsArea.Layout(gtx, g_theme, gridBackground)
				}),
				layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
					editingArea.board.drawingColor = settingsArea.colorPicker.ChosenColor
					return editingArea.Layout(gtx, gridBackground)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
