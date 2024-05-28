package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

var g_theme *material.Theme

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
				layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
					editingArea.board.drawingColor = settingsArea.colorPicker.ChosenColor
					return editingArea.Layout(gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}