package main

import (
	"fmt"
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

type ColorPicker struct {
	size        image.Point
	colors      [9]color.NRGBA
	partLength  int
	chosenColor color.NRGBA
	pickRatio   *float32
	pickerSize  float32
}

func newColorPicker(size image.Point) *ColorPicker {
	cp := &ColorPicker{}
	cp.size = size
	cp.pickerSize = float32(cp.size.Y) * 0.75

	cp.colors = [9]color.NRGBA{
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

	cp.partLength = cp.size.X / (len(cp.colors) - 1)

	return cp
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	cp.size.X = gtx.Constraints.Max.X
	cp.partLength = cp.size.X / (len(cp.colors) - 1)
	cp.pickerSize = float32(cp.size.Y) * 0.4

	cp.HandleInput(gtx)
	cp.Draw(gtx)

	return layout.Dimensions{Size: cp.size}
}

func (cp *ColorPicker) HandleInput(gtx layout.Context) {
	r := image.Rect(9, 9, cp.size.X, cp.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)

	event.Op(gtx.Ops, cp)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cp,
			Kinds:        pointer.Drag,
			ScrollBounds: image.Rect(-10, -10, 10, 10),
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		r := e.Position.X / float32(cp.size.X)
		cp.pickRatio = &r
		*cp.pickRatio = float32(math.Max(math.Min(float64(*cp.pickRatio), 1), 0))

		color := cp.getColorFromPosition(cp.getPickerPositionClamped())
		if color != nil {
			cp.chosenColor = *color
		}
	}

	area.Pop()
}

func (cp *ColorPicker) getPickerPositionClamped() float32 {
	minmax := float64(cp.size.X) / float64(len(cp.colors) - 1) / 2
	pickerPosition := *cp.pickRatio * float32(cp.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cp.size.X) - minmax), minmax))
	return pickerPosition
}

func (cp *ColorPicker) DrawPick(position float32, colorAtPosition color.NRGBA, gtx layout.Context) {
	roundedX := int(math.Round(float64(position)))

	sizeOuter := int(math.Round(float64(cp.pickerSize)))
	sizeInner := int(math.Round(float64(cp.pickerSize) * 0.75))

	fmt.Println(position)

	{
		y0 := cp.size.Y/2 - sizeOuter
		y1 := cp.size.Y/2 + sizeOuter
		r := image.Rect(roundedX-sizeOuter, y0, roundedX+sizeOuter, y1)
		paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}

	{
		y0 := cp.size.Y/2 - sizeInner
		y1 := cp.size.Y/2 + sizeInner
		r := image.Rect(roundedX-sizeInner, y0, roundedX+sizeInner, y1)
		paint.ColorOp{Color: cp.chosenColor}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}
}

func (cp *ColorPicker) Draw(gtx layout.Context) {
	startX := 0
	endX := cp.partLength
	for i := 0; i < len(cp.colors)-1; i++ {
		from := cp.colors[i]
		to := cp.colors[i+1]

		grect := image.Rect(startX, 0, endX, cp.size.Y)
		paint.LinearGradientOp{
			Stop1:  f32.Pt(float32(grect.Min.X), float32(grect.Min.Y)),
			Stop2:  f32.Pt(float32(grect.Max.X), float32(grect.Min.Y)),
			Color1: from,
			Color2: to,
		}.Add(gtx.Ops)
		garea := clip.Rect(grect).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		garea.Pop()

		startX += cp.partLength
		endX += cp.partLength
	}

	if cp.pickRatio != nil {
		p := cp.getPickerPositionClamped()
		col := cp.getColorFromPosition(p)
		if col != nil {
			cp.DrawPick(p, *col, gtx)
		}
	}
}

func (cp *ColorPicker) getColorFromPosition(x float32) *color.NRGBA {
	relX := int(math.Round(float64(x)))
	gradientNumber := relX / cp.partLength
	if gradientNumber >= 0 && gradientNumber < len(cp.colors)-1 {
		col1 := cp.colors[gradientNumber]
		col2 := cp.colors[gradientNumber+1]
		gradientStride := float64(relX%cp.partLength) / float64(cp.partLength)
		colorResult := lerpColor(col1, col2, gradientStride)
		return &colorResult
	}

	return nil
}

func lerpColor(col1 color.NRGBA, col2 color.NRGBA, t float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(int(col1.R) + int(float64(int(col2.R)-int(col1.R))*t)),
		G: uint8(int(col1.G) + int(float64(int(col2.G)-int(col1.G))*t)),
		B: uint8(int(col1.B) + int(float64(int(col2.B)-int(col1.B))*t)),
		A: uint8(int(col1.A) + int(float64(int(col2.A)-int(col1.A))*t)),
	}
}
