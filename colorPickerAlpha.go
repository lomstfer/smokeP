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

type ColorPickerAlpha struct {
	size                     image.Point
	chosenColor              color.NRGBA
	lastColorWithoutAlpha    color.NRGBA
	pickFractionPos          float32
	renderImage              paint.ImageOp
	triggerRenderImageUpdate bool
}

func newColorPickerAlpha(colorWithoutAlpha color.NRGBA, size image.Point) *ColorPickerAlpha {
	cpa := &ColorPickerAlpha{}
	cpa.size = size

	cpa.updateChosenColorFromPickerPos(colorWithoutAlpha)

	return cpa
}

func (cpa *ColorPickerAlpha) Layout(colorWithoutAlpha color.NRGBA, gtx layout.Context) layout.Dimensions {
	cpa.triggerRenderImageUpdate = cpa.triggerRenderImageUpdate || cpa.size != gtx.Constraints.Max
	cpa.size = gtx.Constraints.Max

	cpa.HandleInput(colorWithoutAlpha, gtx)
	cpa.Draw(colorWithoutAlpha, gtx)

	return layout.Dimensions{Size: cpa.size}
}

func (cpa *ColorPickerAlpha) HandleInput(colorWithoutAlpha color.NRGBA, gtx layout.Context) {
	r := image.Rect(0, 0, cpa.size.X, cpa.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)

	event.Op(gtx.Ops, cpa)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cpa,
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

		cpa.pickFractionPos = e.Position.X / float32(cpa.size.X)
	}

	cpa.updateChosenColorFromPickerPos(colorWithoutAlpha)

	area.Pop()
}

func (cpa *ColorPickerAlpha) updateChosenColorFromPickerPos(colorWithoutAlpha color.NRGBA) {
	cpa.pickFractionPos = float32(math.Max(math.Min(float64(cpa.pickFractionPos), 1), 0))

	color := cpa.getColorFromPosition(cpa.getPickerPositionClamped(), colorWithoutAlpha)
	if color != nil {
		cpa.triggerRenderImageUpdate = cpa.triggerRenderImageUpdate || *color != cpa.chosenColor || colorWithoutAlpha != cpa.lastColorWithoutAlpha
		cpa.chosenColor = *color
		cpa.lastColorWithoutAlpha = colorWithoutAlpha
	}
}

func (cpa *ColorPickerAlpha) getPickerPositionClamped() float32 {
	pickerPosition := cpa.pickFractionPos * float32(cpa.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cpa.size.X)), 0))
	return pickerPosition
}

func (cpa *ColorPickerAlpha) Draw(colorWithoutAlpha color.NRGBA, gtx layout.Context) {
	grect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

	if cpa.triggerRenderImageUpdate {
		cpa.triggerRenderImageUpdate = false
		img := image.NewNRGBA(grect)
		for x := 0; x < grect.Dx(); x++ {
			col := cpa.getColorFromPosition(float32(x), colorWithoutAlpha)
			if col != nil {
				for y := 0; y < grect.Dy(); y++ {
					img.SetNRGBA(x, y, *col)
				}
			}
		}
		cpa.renderImage = paint.NewImageOp(img)
	}

	cpa.renderImage.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	p := cpa.getPickerPositionClamped()
	col := cpa.getColorFromPosition(p, colorWithoutAlpha)
	if col != nil {
		drawPicker(f32.Pt(p, float32(cpa.size.Y)/2), *col, gtx)
	}
}

func (cpa *ColorPickerAlpha) getColorFromPosition(x float32, colorWithoutAlpha color.NRGBA) *color.NRGBA {
	relX := int(math.Round(float64(x)))
	colorResult := lerpColor(colorWithoutAlpha, color.NRGBA{0, 0, 0, 0}, float64(relX)/float64(cpa.size.X))
	return &colorResult
}
