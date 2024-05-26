package main

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPickerHue struct {
	size                image.Point
	colors              [7]color.NRGBA
	partLength          int
	chosenColor         color.NRGBA
	pickFractionOfWhole *float32
	pickerSize          float32
}

func newColorPickerHue(size image.Point) *ColorPickerHue {
	cph := &ColorPickerHue{}
	cph.size = size
	cph.pickerSize = float32(cph.size.Y) * 0.75

	cph.colors = [7]color.NRGBA{
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 255, 255, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{255, 0, 0, 255},
	}

	cph.partLength = cph.size.X / (len(cph.colors) - 1)
	cph.updateChosenColorFromPickerPos(0.0)

	return cph
}

func (cph *ColorPickerHue) Layout(gtx layout.Context) layout.Dimensions {
	cph.size = gtx.Constraints.Max
	cph.partLength = cph.size.X / (len(cph.colors) - 1)
	cph.pickerSize = float32(cph.size.Y) * 0.4

	cph.HandleInput(gtx)
	cph.Draw(gtx)

	return layout.Dimensions{Size: cph.size}
}

func (cph *ColorPickerHue) HandleInput(gtx layout.Context) {
	r := image.Rect(0, 0, cph.size.X, cph.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)

	event.Op(gtx.Ops, cph)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cph,
			Kinds:        pointer.Drag | pointer.Press,
			ScrollBounds: image.Rect(-10, -10, 10, 10),
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		if !e.Buttons.Contain(pointer.ButtonPrimary) {
			continue
		}

		r := e.Position.X / float32(cph.size.X)
		cph.updateChosenColorFromPickerPos(r)
	}

	area.Pop()
}

func (cph *ColorPickerHue) updateChosenColorFromPickerPos(fraction float32) {
	cph.pickFractionOfWhole = &fraction
	*cph.pickFractionOfWhole = float32(math.Max(math.Min(float64(*cph.pickFractionOfWhole), 1), 0))

	color := cph.getColorFromPosition(cph.getPickerPositionClamped())
	if color != nil {
		cph.chosenColor = *color
	}
}

func (cph *ColorPickerHue) getPickerPositionClamped() float32 {
	pickerPosition := *cph.pickFractionOfWhole * float32(cph.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cph.size.X)), 0))
	return pickerPosition
}

func (cph *ColorPickerHue) Draw(gtx layout.Context) {
	grect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

    img := image.NewNRGBA(grect)

    for x := 0; x < grect.Dx(); x++ {
		col := cph.getColorFromPosition(float32(x))
		if (col != nil) {
			for y := 0; y < grect.Dy(); y++ {
				img.SetNRGBA(x, y, *col)
			}
		}
    }

	paint.NewImageOp(img).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	if cph.pickFractionOfWhole != nil {
		p := cph.getPickerPositionClamped()
		col := cph.getColorFromPosition(p)
		if col != nil {
			drawPicker(f32.Pt(p, float32(cph.size.Y) / 2), *col, gtx)
		}
	}
}

func (cph *ColorPickerHue) getColorFromPosition(x float32) *color.NRGBA {
	relX := int(math.Round(float64(x)))
	gradientNumber := relX / cph.partLength
	if gradientNumber >= 0 && gradientNumber < len(cph.colors)-1 {
		col1 := cph.colors[gradientNumber]
		col2 := cph.colors[gradientNumber+1]
		gradientStride := float64(relX%cph.partLength) / float64(cph.partLength)
		colorResult := lerpColor(col1, col2, gradientStride)
		return &colorResult
	}

	return nil
}
