package main

import (
	"fmt"
	"image"
	"smokep/colorPicker"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
			d := sa.colorPicker.Layout(g_theme, gtx)
			if sa.colorPicker.PickedNewColor {
				fmt.Println(sa.colorPicker.ChosenColor)
			}
			return d 
		}),
	)
	
	return d
}