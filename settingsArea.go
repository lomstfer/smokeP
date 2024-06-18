package main

import (
	"fmt"
	"image"
	"smokep/colorPicker"
	"smokep/utils"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SettingsArea struct {
	colorPicker                *colorPicker.ColorPicker
	saveButton                 *widget.Clickable
	SaveButtonClicked          chan bool
	loadButton                 *widget.Clickable
	LoadButtonClicked          chan bool
	pixelBoardSizeEditor       widget.Editor
	PixelBoardSizeEditorSubmit chan image.Point
}

func newSettingsArea(pixelBoardSize image.Point) *SettingsArea {
	sa := &SettingsArea{}

	sa.colorPicker = colorPicker.NewColorPicker(image.Pt(1, 1))
	sa.saveButton = &widget.Clickable{}
	sa.loadButton = &widget.Clickable{}
	sa.SaveButtonClicked = make(chan bool)
	sa.LoadButtonClicked = make(chan bool)
	sa.pixelBoardSizeEditor = widget.Editor{Submit: true, ReadOnly: false}
	sa.pixelBoardSizeEditor.SetText(fmt.Sprintf("%dx%d", pixelBoardSize.X, pixelBoardSize.Y))
	sa.PixelBoardSizeEditorSubmit = make(chan image.Point)

	return sa
}

func (sa *SettingsArea) Update(gtx layout.Context, pixelBoardSize image.Point) {
	utils.ConsumePressAndFocusSelf(sa, gtx)

	sa.colorPicker.Update(gtx)

	if !gtx.Focused(&sa.pixelBoardSizeEditor) {
		sa.pixelBoardSizeEditor.SetText(fmt.Sprintf("%dx%d", pixelBoardSize.X, pixelBoardSize.Y))
	}

	for {
		ev, ok := sa.pixelBoardSizeEditor.Update(gtx)
		if !ok {
			break
		}
		_, ok = ev.(widget.SubmitEvent)
		if !ok {
			continue
		}
		gtx.Execute(key.FocusCmd{Tag: nil})

		{
			input := sa.pixelBoardSizeEditor.Text()
			var width, height int
			_, err := fmt.Sscanf(input, "%dx%d", &width, &height)
			if err == nil {
				sa.PixelBoardSizeEditorSubmit <- image.Pt(width, height)
			}
		}
	}
}

func (sa *SettingsArea) Layout(gtx layout.Context, theme *material.Theme, gridBg *utils.GridBackground) layout.Dimensions {
	{
		area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Push(gtx.Ops)
		event.Op(gtx.Ops, sa)
		area.Pop()
	}

	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			height := 0
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sa.saveButton.Clicked(gtx) {
						sa.SaveButtonClicked <- true
					}
					d := material.Button(theme, sa.saveButton, "Save").Layout(gtx)
					height += d.Size.Y
					return d
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sa.loadButton.Clicked(gtx) {
						sa.LoadButtonClicked <- true
					}
					d := material.Button(theme, sa.loadButton, "Load").Layout(gtx)
					height += d.Size.Y
					return d
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := material.Editor(theme, &sa.pixelBoardSizeEditor, "").Layout(gtx)
					height += d.Size.Y

					return d
				}),
			)

			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, height)}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			gridBg.Draw(gtx.Ops, r)
			area := clip.Rect(r).Push(gtx.Ops)
			paint.ColorOp{Color: sa.colorPicker.ChosenColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			area.Pop()
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			area := clip.Rect(r).Push(gtx.Ops)
			area.Pop()
			d := sa.colorPicker.Layout(gtx, g_theme, gridBg)
			return d
		}),
	)

	return d
}
