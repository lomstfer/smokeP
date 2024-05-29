package main

import (
	"fmt"
	"image"
	"smokep/colorPicker"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type SettingsArea struct {
	colorPicker *colorPicker.ColorPicker
}

func newSettingsArea() *SettingsArea {
	sa := &SettingsArea{}

	sa.colorPicker = colorPicker.NewColorPicker(image.Pt(100, 50))
	return sa
}

func (sa *SettingsArea) Layout(gtx layout.Context) layout.Dimensions {
	{
		area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Push(gtx.Ops)
		paint.ColorOp{Color: g_theme.Fg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		area.Pop()
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
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
			return sa.colorPicker.Layout(g_theme, gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			area := clip.Rect(r).Push(gtx.Ops)
			area.Pop()

			t := material.H2(g_theme, fmt.Sprintf("%v", sa.colorPicker.ChosenColor))
			t.Color = sa.colorPicker.ChosenColor
			t.Alignment = text.Middle
			return t.Layout(gtx)
		}),
	)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
