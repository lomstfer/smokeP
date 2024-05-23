package main

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPicker struct {
	size     image.Point
	hue      *ColorPickerHue
	valSat *ColorPickerValueSat
	chosenColor color.NRGBA
}

func newColorPicker(size image.Point) *ColorPicker {
	cp := &ColorPicker{}
	cp.hue = newColorPickerHue(size)
	cp.valSat = newColorPickerValueSat(cp.hue.chosenColor, size)

	return cp
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return cp.hue.Layout(gtx)
		}),
		layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
			d := cp.valSat.Layout(cp.hue.chosenColor, gtx)
			cp.chosenColor = cp.valSat.chosenColor
			return d
		}),
	)

	return d
}

func lerpColor(col1 color.NRGBA, col2 color.NRGBA, t float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(int(col1.R) + int(float64(int(col2.R)-int(col1.R))*t)),
		G: uint8(int(col1.G) + int(float64(int(col2.G)-int(col1.G))*t)),
		B: uint8(int(col1.B) + int(float64(int(col2.B)-int(col1.B))*t)),
		A: uint8(int(col1.A) + int(float64(int(col2.A)-int(col1.A))*t)),
	}
}

func drawPicker(position f32.Point, colorAtPosition color.NRGBA, gtx layout.Context) {
	const pickerSize = 10

	roundedPos := position.Round()

	sizeOuter := int(math.Round(float64(pickerSize)))
	sizeInner := int(math.Round(float64(pickerSize) * 0.75))

	{
		y0 := roundedPos.Y - sizeOuter
		y1 := roundedPos.Y + sizeOuter
		r := image.Rect(roundedPos.X-sizeOuter, y0, roundedPos.X+sizeOuter, y1)
		paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}

	{
		y0 := roundedPos.Y - sizeInner
		y1 := roundedPos.Y + sizeInner
		r := image.Rect(roundedPos.X-sizeInner, y0, roundedPos.X+sizeInner, y1)
		paint.ColorOp{Color: colorAtPosition}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}
}
