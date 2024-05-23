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
	colors              [9]color.NRGBA
	partLength          int
	chosenColor         color.NRGBA
	pickFractionOfWhole *float32
	pickerSize          float32
}

func newColorPickerHue(size image.Point) *ColorPickerHue {
	cph := &ColorPickerHue{}
	cph.size = size
	cph.pickerSize = float32(cph.size.Y) * 0.75

	cph.colors = [9]color.NRGBA{
		{255, 0, 0, 255},
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 255, 255, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{255, 0, 0, 255},
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
	minmax := float64(cph.size.X) / float64(len(cph.colors)-1) / 2
	pickerPosition := *cph.pickFractionOfWhole * float32(cph.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cph.size.X)-minmax), minmax))
	return pickerPosition
}

func (cph *ColorPickerHue) Draw(gtx layout.Context) {
	startX := 0
	endX := cph.partLength
	for i := 0; i < len(cph.colors)-1; i++ {
		from := cph.colors[i]
		to := cph.colors[i+1]

		grect := image.Rect(startX, 0, endX, cph.size.Y)
		paint.LinearGradientOp{
			Stop1:  f32.Pt(float32(grect.Min.X), float32(grect.Min.Y)),
			Stop2:  f32.Pt(float32(grect.Max.X), float32(grect.Min.Y)),
			Color1: from,
			Color2: to,
		}.Add(gtx.Ops)
		garea := clip.Rect(grect).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		garea.Pop()

		startX += cph.partLength
		endX += cph.partLength
	}

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
