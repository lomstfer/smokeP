package colorPicker

import (
	"image"
	"image/color"
	"math"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type ColorPicker struct {
	hue         *ColorPickerHue
	valSat      *ColorPickerValueSat
	alpha       *ColorPickerAlpha
	ChosenColor color.NRGBA
}

func NewColorPicker(size image.Point) *ColorPicker {
	cp := &ColorPicker{}
	cp.hue = newColorPickerHue(size)
	cp.valSat = newColorPickerValueSat(cp.hue.chosenColor, size)
	cp.ChosenColor = cp.valSat.chosenColor
	cp.alpha = newColorPickerAlpha(cp.ChosenColor, size)

	return cp
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.hue.Layout(gtx)
			return d
		}),
		layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
			d := cp.valSat.Layout(cp.hue.chosenColor, gtx)
			return d
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			cp.ChosenColor = cp.valSat.chosenColor
			d := cp.alpha.Layout(cp.ChosenColor, gtx)
			cp.ChosenColor.A = cp.alpha.chosenColor.A
			return d
		}),
	)

	return d
}

func lerpColor(col1 color.NRGBA, col2 color.NRGBA, t float64) color.NRGBA {
	t = utils.Clamp(t, 0, 1)
	return color.NRGBA{
		R: uint8(int(col1.R) + int(float64(int(col2.R)-int(col1.R))*t)),
		G: uint8(int(col1.G) + int(float64(int(col2.G)-int(col1.G))*t)),
		B: uint8(int(col1.B) + int(float64(int(col2.B)-int(col1.B))*t)),
		A: uint8(int(col1.A) + int(float64(int(col2.A)-int(col1.A))*t)),
	}
}

func drawPicker(position f32.Point, colorAtPosition color.NRGBA, gtx layout.Context) {
	pickerSize := unit.Dp(15)

	roundedPos := position.Round()

	sizeOuterOuter := int(math.Round(float64(pickerSize)))
	sizeOuter := int(math.Round(float64(pickerSize)) * 0.8)
	sizeInner := int(math.Round(float64(sizeOuter) * 0.8))

	// lightness := 0.299*float64(colorAtPosition.R)/255 + 0.587*float64(colorAtPosition.G)/255 + 0.114*float64(colorAtPosition.B)/255
	// var outlineColor color.NRGBA
	// if lightness < 0.5 {
	// 	outlineColor = color.NRGBA{255, 255, 255, 255}
	// } else {
	// 	outlineColor = color.NRGBA{0, 0, 0, 255}
	// }

	{
		y0 := roundedPos.Y - sizeOuterOuter
		y1 := roundedPos.Y + sizeOuterOuter
		r := image.Rect(roundedPos.X-sizeOuterOuter, y0, roundedPos.X+sizeOuterOuter, y1)
		paint.ColorOp{Color: color.NRGBA{0,0,0,255}}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}

	{
		y0 := roundedPos.Y - sizeOuter
		y1 := roundedPos.Y + sizeOuter
		r := image.Rect(roundedPos.X-sizeOuter, y0, roundedPos.X+sizeOuter, y1)
		paint.ColorOp{Color: color.NRGBA{255,255,255,255}}.Add(gtx.Ops)
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
