package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type SettingsArea struct {
	colorPicker *ColorPicker
}

func newSettingsArea() *SettingsArea {
	sa := &SettingsArea{}

	sa.colorPicker = newColorPicker(image.Pt(100, 50))

	return sa
}

func (sa *SettingsArea) Layout(gtx layout.Context) layout.Dimensions {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			title := material.H1(g_theme, "settings")
			title.Color = color.NRGBA{0, 0, 0, 255}
			title.Alignment = text.Middle

			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			area := clip.Rect(r).Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{100, 100, 100, 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			area.Pop()

			return title.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			area := clip.Rect(r).Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{255, 100, 100, 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			area.Pop()
			return sa.colorPicker.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			area := clip.Rect(r).Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{100, 100, 255, 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			area.Pop()

			t := material.H2(g_theme, fmt.Sprintf("%v", sa.colorPicker.chosenColor))
			t.Color = sa.colorPicker.chosenColor
			t.Alignment = text.Middle
			return t.Layout(gtx)
		}),
	)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
