package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPicker struct {
	position f32.Point
	size     image.Point
	colors   [9]color.NRGBA
	partLength int
	chosenColor color.NRGBA
}

func newColorPicker(position f32.Point, size image.Point) *ColorPicker {
	cp := &ColorPicker{}
	cp.position = position
	cp.size = size

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
	cp.partLength = cp.size.X / (len(cp.colors) - 1)

	cp.HandleInput(gtx)
	cp.Draw(gtx)

	return layout.Dimensions{Size: cp.size}
}

func (cp *ColorPicker) HandleInput(gtx layout.Context) {
	roundPos := cp.position.Round()

	r := image.Rect(roundPos.X, roundPos.Y, roundPos.X+cp.size.X, roundPos.Y+cp.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)

	event.Op(gtx.Ops, cp)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cp,
			Kinds:        pointer.Move,
			ScrollBounds: image.Rect(-10, -10, 10, 10),
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		relX := e.Position.Sub(cp.position).Round().X
		gradientNumber := relX / cp.partLength
		if (gradientNumber < len(cp.colors) - 1) {
			col1 := cp.colors[gradientNumber]
			col2 := cp.colors[gradientNumber + 1]
			gradientStride := float64(relX % cp.partLength) / float64(cp.partLength)
			
			cp.chosenColor = lerpColor(col1, col2, gradientStride)
		}
	}

	area.Pop()
}

func (cp *ColorPicker) Draw(gtx layout.Context) {
	roundPos := cp.position.Round()

	startX := roundPos.X
	endX := roundPos.X + cp.partLength
	for i := 0; i < len(cp.colors)-1; i++ {
		from := cp.colors[i]
		to := cp.colors[i+1]

		grect := image.Rect(startX, roundPos.Y, endX, roundPos.Y+cp.size.Y)
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
}

func lerpColor(col1 color.NRGBA, col2 color.NRGBA, t float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(int(col1.R) + int(float64(int(col2.R) - int(col1.R)) * t)),
		G: uint8(int(col1.G) + int(float64(int(col2.G) - int(col1.G)) * t)),
		B: uint8(int(col1.B) + int(float64(int(col2.B) - int(col1.B)) * t)),
		A: uint8(int(col1.A) + int(float64(int(col2.A) - int(col1.A)) * t)),
	}
}