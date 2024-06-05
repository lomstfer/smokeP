package colorPicker

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ColorPicker struct {
	hue             *ColorPickerHue
	valSat          *ColorPickerValueSat
	alpha           *ColorPickerAlpha
	ChosenColor     color.NRGBA
	rgbaEditor      widget.Editor
	hexEditor       widget.Editor
	rgbaEditorFocus bool
	hexEditorFocus  bool
	PickedNewColor  bool
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
	cp.PickedNewColor = false

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.hue.Layout(gtx)
			if cp.hue.pickedNewColor {
				cp.updateColors()
			}
			return d
		}),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			d := cp.valSat.Layout(cp.hue.chosenColor, gtx)
			if cp.valSat.pickedNewColor {
				cp.updateColors()
			}

			return d
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.alpha.Layout(cp.ChosenColor, gtx)
			if cp.alpha.pickedNewColor {
				cp.updateColors()
			}

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
						e, ok := ev.(widget.SubmitEvent)
						if !ok {
							continue
						}
						fmt.Println(e.Text)
						gtx.Execute(key.FocusCmd{Tag: nil})
						cp.rgbaEditor.SetText(fmt.Sprintf("rgba(%v, %v, %v, %v)", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
						cp.setPickersToColor(color.NRGBA{100, 200, 200, 255})
					}

					return material.Editor(theme, &cp.rgbaEditor, "").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					for {
						ev, ok := cp.hexEditor.Update(gtx)
						if !ok {
							break
						}
						e, ok := ev.(widget.SubmitEvent)
						if !ok {
							continue
						}
						fmt.Println(e.Text)
						gtx.Execute(key.FocusCmd{Tag: nil})
						cp.hexEditor.SetText(fmt.Sprintf("#%02x%02x%02x%02x", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
					}

					return material.Editor(theme, &cp.hexEditor, "").Layout(gtx)
				}),
			)
		}),
	)
	
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.hue.Layout(gtx)
			if cp.hue.pickedNewColor {
				cp.updateColors()
			}
			return d
		}),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			d := cp.valSat.Layout(cp.hue.chosenColor, gtx)
			if cp.valSat.pickedNewColor {
				cp.updateColors()
			}

			return d
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := cp.alpha.Layout(cp.ChosenColor, gtx)
			if cp.alpha.pickedNewColor {
				cp.updateColors()
			}

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
						e, ok := ev.(widget.SubmitEvent)
						if !ok {
							continue
						}
						fmt.Println(e.Text)
						gtx.Execute(key.FocusCmd{Tag: nil})
						cp.rgbaEditor.SetText(fmt.Sprintf("rgba(%v, %v, %v, %v)", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
						cp.setPickersToColor(color.NRGBA{100, 200, 200, 255})
					}

					return material.Editor(theme, &cp.rgbaEditor, "").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					for {
						ev, ok := cp.hexEditor.Update(gtx)
						if !ok {
							break
						}
						e, ok := ev.(widget.SubmitEvent)
						if !ok {
							continue
						}
						fmt.Println(e.Text)
						gtx.Execute(key.FocusCmd{Tag: nil})
						cp.hexEditor.SetText(fmt.Sprintf("#%02x%02x%02x%02x", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
					}

					return material.Editor(theme, &cp.hexEditor, "").Layout(gtx)
				}),
			)
		}),
	)

	cp.rgbaEditorFocus = gtx.Focused(&cp.rgbaEditor)
	cp.hexEditorFocus = gtx.Focused(&cp.hexEditor)

	return d
}

func (cp *ColorPicker) updateColors() {
	cp.PickedNewColor = true
	cp.valSat.updateChosenColor(cp.hue.chosenColor)
	cp.alpha.updateChosenColor(cp.valSat.chosenColor)
	cp.ChosenColor = cp.valSat.chosenColor
	cp.ChosenColor.A = cp.alpha.chosenAlpha

	if !cp.rgbaEditorFocus {
		cp.rgbaEditor.SetText(fmt.Sprintf("rgba(%v, %v, %v, %v)", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
	}
	if !cp.hexEditorFocus {
		cp.hexEditor.SetText(fmt.Sprintf("#%02x%02x%02x%02x", cp.ChosenColor.R, cp.ChosenColor.G, cp.ChosenColor.B, cp.ChosenColor.A))
	}
}

func (cp *ColorPicker) setPickersToColor(newColor color.NRGBA) {
	cp.hue.updateColor(newColor)
	cp.valSat.updateColor(newColor, cp.hue.chosenColor)
	cp.alpha.updateColor(cp.valSat.chosenColor, newColor.A)
	cp.ChosenColor = cp.valSat.chosenColor
	cp.ChosenColor.A = cp.alpha.chosenAlpha
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
	DrawCircleOutline(gtx, position, 10, 2, color.NRGBA{255, 255, 255, 255})
	DrawCircleOutline(gtx, position, 8, 2, color.NRGBA{0, 0, 0, 255})
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
