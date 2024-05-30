package colorPicker

import (
	"image"
	"image/color"
	"math"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPickerAlpha struct {
	size                     image.Point
	chosenColor              color.NRGBA
	lastopaqueChosenColor    color.NRGBA
	pickFractionPos          float32
	renderImage              paint.ImageOp
	triggerRenderImageUpdate bool
	background paint.ImageOp
}

func newColorPickerAlpha(opaqueChosenColor color.NRGBA, size image.Point) *ColorPickerAlpha {
	cpa := &ColorPickerAlpha{}
	cpa.size = size

	cpa.updateChosenColorFromPickerPos(opaqueChosenColor)
	cpa.background = paint.NewImageOp(utils.GenerateGridImage(50, 25, color.NRGBA{200, 200, 200, 255}, color.NRGBA{100, 100, 100, 255}))
	cpa.background.Filter = paint.FilterNearest

	return cpa
}

func (cpa *ColorPickerAlpha) Layout(opaqueChosenColor color.NRGBA, gtx layout.Context) layout.Dimensions {
	cpa.triggerRenderImageUpdate = cpa.triggerRenderImageUpdate || cpa.size != gtx.Constraints.Max
	cpa.size = gtx.Constraints.Max

	cpa.HandleInput(opaqueChosenColor, gtx)
	cpa.Draw(opaqueChosenColor, gtx)

	return layout.Dimensions{Size: cpa.size}
}

func (cpa *ColorPickerAlpha) HandleInput(opaqueChosenColor color.NRGBA, gtx layout.Context) {
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

	cpa.updateChosenColorFromPickerPos(opaqueChosenColor)

	area.Pop()
}

func (cpa *ColorPickerAlpha) updateChosenColorFromPickerPos(opaqueChosenColor color.NRGBA) {
	cpa.pickFractionPos = float32(math.Max(math.Min(float64(cpa.pickFractionPos), 1), 0))

	color := cpa.getColorFromPosition(cpa.getPickerPositionClamped(), opaqueChosenColor)
	if color != nil {
		cpa.triggerRenderImageUpdate = cpa.triggerRenderImageUpdate || *color != cpa.chosenColor || opaqueChosenColor != cpa.lastopaqueChosenColor
		cpa.chosenColor = *color
		cpa.lastopaqueChosenColor = opaqueChosenColor
	}
}

func (cpa *ColorPickerAlpha) getPickerPositionClamped() float32 {
	pickerPosition := cpa.pickFractionPos * float32(cpa.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cpa.size.X)), 0))
	return pickerPosition
}

func (cpa *ColorPickerAlpha) Draw(opaqueChosenColor color.NRGBA, gtx layout.Context) {
	cpa.drawBackground(gtx.Ops)
	
	grect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

	if cpa.triggerRenderImageUpdate {
		cpa.triggerRenderImageUpdate = false
		img := image.NewNRGBA(grect)
		for x := 0; x < grect.Dx(); x++ {
			col := cpa.getColorFromPosition(float32(x), opaqueChosenColor)
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
	col := cpa.getColorFromPosition(p, opaqueChosenColor)
	if col != nil {
		drawPicker(f32.Pt(p, float32(cpa.size.Y)/2), *col, gtx)
	}
}

func (cpa *ColorPickerAlpha) getColorFromPosition(x float32, opaqueChosenColor color.NRGBA) *color.NRGBA {
	colorResult := lerpColor(opaqueChosenColor, color.NRGBA{0, 0, 0, 0}, float64(x) / float64(cpa.size.X))
	return &colorResult
}

func (cpa *ColorPickerAlpha) drawBackground(ops *op.Ops) {
	area := clip.Rect(image.Rect(0, 0, cpa.size.X, cpa.size.Y)).Push(ops)
	
	cpa.background.Add(ops)
	
	scale := max(float32(cpa.size.X) / float32(cpa.background.Size().X), float32(cpa.size.Y) / float32(cpa.background.Size().Y))
	tStack := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(ops)
	paint.PaintOp{}.Add(ops)
	tStack.Pop()

	area.Pop()
}