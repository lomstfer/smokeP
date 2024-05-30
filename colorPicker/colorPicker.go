package colorPicker

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ColorPicker struct {
	hue         *ColorPickerHue
	valSat      *ColorPickerValueSat
	alpha       *ColorPickerAlpha
	ChosenColor color.NRGBA
	rgbaEditor  widget.Editor
	hexEditor   widget.Editor
}

func NewColorPicker(size image.Point) *ColorPicker {
	cp := &ColorPicker{}
	cp.hue = newColorPickerHue(size)
	cp.valSat = newColorPickerValueSat(cp.hue.chosenColor, size)
	cp.ChosenColor = cp.valSat.chosenColor
	cp.alpha = newColorPickerAlpha(cp.ChosenColor, size)

	cp.rgbaEditor = widget.Editor{Submit: true, ReadOnly: true}
	cp.hexEditor = widget.Editor{Submit: true, ReadOnly: true}

	return cp
}

func (cp *ColorPicker) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.hue.Layout(gtx)
			return d
		}),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			d := cp.valSat.Layout(cp.hue.chosenColor, gtx)

			cp.ChosenColor = cp.valSat.chosenColor
			cp.ChosenColor.A = cp.alpha.chosenColor.A
			cp.updateChosenColorText()

			return d
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			cp.ChosenColor = cp.valSat.chosenColor
			d := cp.alpha.Layout(cp.ChosenColor, gtx)
			cp.ChosenColor.A = cp.alpha.chosenColor.A

			cp.updateChosenColorText()

			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Editor(theme, &cp.rgbaEditor, "").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Editor(theme, &cp.hexEditor, "").Layout(gtx)
				}),
			)
		}),
	)

	return d
}

func (cp *ColorPicker) updateChosenColorText() {
	cp.rgbaEditor.SetText(fmt.Sprintf("rgba(%v, %v, %v, %v)", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
	cp.hexEditor.SetText(fmt.Sprintf("#%02x%02x%02x%02x", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
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
	DrawCircleOutline(gtx, position, float32(unit.Dp(6)), 2, color.NRGBA{255, 255, 255, 255})
	DrawCircleOutline(gtx, position, float32(unit.Dp(4)), 2, color.NRGBA{0, 0, 0, 255})

	// pickerSize := unit.Dp(15)

	// roundedPos := position.Round()

	// sizeOuterOuter := int(math.Round(float64(pickerSize)))
	// sizeOuter := int(math.Round(float64(pickerSize)) * 0.8)
	// sizeInner := int(math.Round(float64(sizeOuter) * 0.8))

	// // lightness := 0.299*float64(colorAtPosition.R)/255 + 0.587*float64(colorAtPosition.G)/255 + 0.114*float64(colorAtPosition.B)/255
	// // var outlineColor color.NRGBA
	// // if lightness < 0.5 {
	// // 	outlineColor = color.NRGBA{255, 255, 255, 255}
	// // } else {
	// // 	outlineColor = color.NRGBA{0, 0, 0, 255}
	// // }

	// {
	// 	y0 := roundedPos.Y - sizeOuterOuter
	// 	y1 := roundedPos.Y + sizeOuterOuter
	// 	r := image.Rect(roundedPos.X-sizeOuterOuter, y0, roundedPos.X+sizeOuterOuter, y1)
	// 	paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
	// 	a := clip.Ellipse(r).Push(gtx.Ops)
	// 	paint.PaintOp{}.Add(gtx.Ops)
	// 	a.Pop()
	// }

	// {
	// 	y0 := roundedPos.Y - sizeOuter
	// 	y1 := roundedPos.Y + sizeOuter
	// 	r := image.Rect(roundedPos.X-sizeOuter, y0, roundedPos.X+sizeOuter, y1)
	// 	paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
	// 	a := clip.Ellipse(r).Push(gtx.Ops)
	// 	paint.PaintOp{}.Add(gtx.Ops)
	// 	a.Pop()
	// }

	// {
	// 	y0 := roundedPos.Y - sizeInner
	// 	y1 := roundedPos.Y + sizeInner
	// 	r := image.Rect(roundedPos.X-sizeInner, y0, roundedPos.X+sizeInner, y1)
	// 	paint.ColorOp{Color: colorAtPosition}.Add(gtx.Ops)
	// 	a := clip.Ellipse(r).Push(gtx.Ops)
	// 	paint.PaintOp{}.Add(gtx.Ops)
	// 	a.Pop()
	// }
}

func DrawCircleOutline(gtx layout.Context, center f32.Point, radius float32, strokeWidth float32, col color.NRGBA) {
	circle := image.Rect(
		int(math.Round(float64(center.X-radius-strokeWidth/2.0))), int(math.Round(float64(center.Y-radius-strokeWidth/2.0))),
		int(math.Round(float64(center.X+radius+strokeWidth/2.0))), int(math.Round(float64(center.Y+radius+strokeWidth/2.0))),
	).Canon()

	// roundedRadius := int(math.Round(float64(radius)))
	s := clip.Stroke{
		Path:  clip.RRect{Rect: circle, SE: circle.Dx()/2.0, SW: circle.Dx()/2.0, NW: circle.Dx()/2.0, NE: circle.Dx()/2.0}.Path(gtx.Ops),
		Width: strokeWidth,
	}.Op().Push(gtx.Ops)
	
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	s.Pop()
}
