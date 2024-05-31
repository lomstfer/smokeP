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

	cp.rgbaEditor = widget.Editor{Submit: true, ReadOnly: false}
	cp.hexEditor = widget.Editor{Submit: true, ReadOnly: true}
	cp.rgbaEditor.SetText("rgba(255, 255, 255, 255)")
	cp.hexEditor.SetText("#ffffffff")

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
					for {
						ev, ok := cp.rgbaEditor.Update(gtx)
						if !ok {
							break
						}
						fmt.Println(gtx.Focused(&cp.rgbaEditor))
						
						fmt.Println("hej", ev)
						e, ok := ev.(widget.SubmitEvent)
						if !ok {
							continue
						}
						fmt.Println(e.Text)
					}
					
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
	// cp.rgbaEditor.SetText(fmt.Sprintf("rgba(%v, %v, %v, %v)", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
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
}

func DrawCircleOutline(gtx layout.Context, center f32.Point, radius float32, strokeWidth float32, col color.NRGBA) {
	circle := image.Rect(
		int(math.Round(float64(center.X-radius-strokeWidth/2.0))), int(math.Round(float64(center.Y-radius-strokeWidth/2.0))),
		int(math.Round(float64(center.X+radius+strokeWidth/2.0))), int(math.Round(float64(center.Y+radius+strokeWidth/2.0))),
	).Canon()

	realRadius := circle.Dx() / 2.0
	s := clip.Stroke{
		Path:  clip.RRect{Rect: circle, SE: realRadius, SW: realRadius, NW: realRadius, NE: realRadius}.Path(gtx.Ops),
		Width: strokeWidth,
	}.Op().Push(gtx.Ops)

	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	s.Pop()
}
