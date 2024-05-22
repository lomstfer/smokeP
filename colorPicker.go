package main

import (
	"fmt"
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
	position     f32.Point
	size         image.Point
	colors       [9]color.NRGBA
	partLength   int
	chosenColor  color.NRGBA
	pickPosition *f32.Point
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

		cp.pickPosition = &e.Position

		color := cp.getColorFromPosition(*cp.pickPosition)
		if color != nil {
			cp.chosenColor = *color
		}
	}

	area.Pop()
}

func (cp *ColorPicker) DrawPick(mousePos f32.Point, colorAtPosition color.NRGBA, gtx layout.Context) {
	const size = 10
	const sizeOuter = 13

	roundedMPos := mousePos.Round()
	fmt.Println(roundedMPos)
	
	{
		y0 := cp.position.Round().Y + cp.size.Y/2 - sizeOuter
		y1 := cp.position.Round().Y + cp.size.Y/2 + sizeOuter
		r := image.Rect(roundedMPos.X-sizeOuter, y0, roundedMPos.X+sizeOuter, y1)
		paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}
	
	{
		y0 := cp.position.Round().Y + cp.size.Y/2 - size
		y1 := cp.position.Round().Y + cp.size.Y/2 + size
		r := image.Rect(roundedMPos.X-size, y0, roundedMPos.X+size, y1)
		paint.ColorOp{Color: cp.chosenColor}.Add(gtx.Ops)
		a := clip.Ellipse(r).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		a.Pop()
	}
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

	if cp.pickPosition != nil {
		col := cp.getColorFromPosition(*cp.pickPosition)
		if col != nil {
			cp.DrawPick(*cp.pickPosition, *col, gtx)
		}
	}
}

func (cp *ColorPicker) getColorFromPosition(pos f32.Point) *color.NRGBA {
	relX := pos.Sub(cp.position).Round().X
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
