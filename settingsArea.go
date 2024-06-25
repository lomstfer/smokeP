package main

import (
	"image"
	"smokep/colorPicker"
	"smokep/utils"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SettingsArea struct {
	colorPicker          *colorPicker.ColorPicker
	saveButton           *widget.Clickable
	SaveButtonClicked    chan bool
	loadButton           *widget.Clickable
	LoadButtonClicked    chan bool
	PixelBoardSizeEditor *BoardSizeEditor
	ColorSaver           *ColorSaver
}

func newSettingsArea(pixelBoardSize image.Point) *SettingsArea {
	sa := &SettingsArea{}

	sa.ColorSaver = NewColorSaver()

	sa.colorPicker = colorPicker.NewColorPicker(image.Pt(1, 1))
	sa.saveButton = &widget.Clickable{}
	sa.loadButton = &widget.Clickable{}
	sa.SaveButtonClicked = make(chan bool)
	sa.LoadButtonClicked = make(chan bool)
	sa.PixelBoardSizeEditor = NewBoardSizeEditor(pixelBoardSize)

	return sa
}

func (sa *SettingsArea) Update(gtx layout.Context, pixelBoardSize image.Point) {
	utils.ConsumePressAndFocusSelf(sa, gtx)

	sa.ColorSaver.Update(gtx, sa.colorPicker.ChosenColor)
	if sa.ColorSaver.justClickedColor != nil {
		sa.colorPicker.SetChosenColor(*sa.ColorSaver.justClickedColor)
	}

	sa.colorPicker.Update(gtx)
	sa.PixelBoardSizeEditor.Update(gtx, pixelBoardSize)
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
						gtx.Execute(key.FocusCmd{Tag: nil})
					}
					d := material.Button(theme, sa.saveButton, "Save").Layout(gtx)
					height += d.Size.Y
					return d
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sa.loadButton.Clicked(gtx) {
						sa.LoadButtonClicked <- true
						gtx.Execute(key.FocusCmd{Tag: nil})
					}
					d := material.Button(theme, sa.loadButton, "Load").Layout(gtx)
					height += d.Size.Y
					return d
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := sa.PixelBoardSizeEditor.Layout(gtx, theme)
					height += d.Size.Y

					return d
				}),
			)

			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, height)}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := sa.colorPicker.Layout(gtx, g_theme, gridBg)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			d := sa.ColorSaver.Layout(gtx, g_theme, gridBg)
			return d
		}),
	)

	return d
}
