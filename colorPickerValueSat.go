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

type ColorPickerValueSat struct {
	size            image.Point
	chosenColor     color.NRGBA
	pickFractionPos *f32.Point
	pickerSize      float32
}

func newColorPickerValueSat(hueColor color.NRGBA, size image.Point) *ColorPickerValueSat {
	cpvs := &ColorPickerValueSat{}
	cpvs.size = size
	cpvs.pickerSize = 10

	cpvs.pickFractionPos = &f32.Point{0, 0}

	cpvs.updateChosenColorFromPickerPos(hueColor)

	return cpvs
}

func (cpvs *ColorPickerValueSat) Layout(hueColor color.NRGBA, gtx layout.Context) layout.Dimensions {
	cpvs.size = gtx.Constraints.Max

	cpvs.HandleInput(hueColor, gtx)
	cpvs.Draw(hueColor, gtx)

	return layout.Dimensions{Size: cpvs.size}
}

func (cpvs *ColorPickerValueSat) HandleInput(hueColor color.NRGBA, gtx layout.Context) {
	r := image.Rect(0, 0, cpvs.size.X, cpvs.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)

	event.Op(gtx.Ops, cpvs)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cpvs,
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

		cpvs.pickFractionPos = &f32.Point{X: e.Position.X / float32(cpvs.size.X), Y: e.Position.Y / float32(cpvs.size.Y)}
	}

	cpvs.updateChosenColorFromPickerPos(hueColor)

	area.Pop()
}

func (cpvs *ColorPickerValueSat) updateChosenColorFromPickerPos(hueColor color.NRGBA) {
	*cpvs.pickFractionPos = f32.Pt(
		float32(math.Max(math.Min(float64(cpvs.pickFractionPos.X), 1), 0)),
		float32(math.Max(math.Min(float64(cpvs.pickFractionPos.Y), 1), 0)))

	color := cpvs.getColorFromPosition(cpvs.getPickerPositionClamped(), hueColor)
	if color != nil {
		cpvs.chosenColor = *color
	}
}

func (cpvs *ColorPickerValueSat) getPickerPositionClamped() f32.Point {
	pickerPosition := f32.Point{
		X: cpvs.pickFractionPos.X * float32(cpvs.size.X),
		Y: cpvs.pickFractionPos.Y * float32(cpvs.size.Y)}
	pickerPosition = f32.Pt(
		float32(math.Max(math.Min(float64(pickerPosition.X), float64(cpvs.size.X)), 0)),
		float32(math.Max(math.Min(float64(pickerPosition.Y), float64(cpvs.size.Y)), 0)))
	return pickerPosition
}

func (cpvs *ColorPickerValueSat) Draw(hueColor color.NRGBA, gtx layout.Context) {
	grect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

	img := image.NewNRGBA(grect)

	for x := 0; x < grect.Dx(); x++ {
		for y := 0; y < grect.Dy(); y++ {
			col := *cpvs.getColorFromPosition(f32.Pt(float32(x), float32(y)), hueColor)
			img.SetNRGBA(x, y, col)
		}
	}

	paint.NewImageOp(img).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	if cpvs.pickFractionPos != nil {
		p := cpvs.getPickerPositionClamped()
		col := cpvs.getColorFromPosition(p, hueColor)
		if col != nil {
			drawPicker(p, *col, gtx)
		}
	}
}

func (cpvs *ColorPickerValueSat) getColorFromPosition(pos f32.Point, hueColor color.NRGBA) *color.NRGBA {
	rel := f32.Pt(pos.X/float32(cpvs.size.X), pos.Y/float32(cpvs.size.Y))
	rgb := lerpColor(color.NRGBA{255, 255, 255, 255}, hueColor, float64(rel.X))
	value := float32(lerpColor(color.NRGBA{0, 0, 0, 255}, color.NRGBA{0, 0, 0, 0}, float64(rel.Y)).A) / 255
	rgb = color.NRGBA{R: uint8(float32(rgb.R) * value), G: uint8(float32(rgb.G) * value), B: uint8(float32(rgb.B) * value), A: rgb.A}

	return &rgb
}
