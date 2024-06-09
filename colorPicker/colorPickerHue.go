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
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPickerHue struct {
	size                     image.Point
	colors                   [7]color.NRGBA
	partLength               float32
	chosenColor              color.NRGBA
	pickFractionPos          float32
	renderImage              paint.ImageOp
	triggerRenderImageUpdate bool
	pickedNewColor           bool
}

func newColorPickerHue(size image.Point) *ColorPickerHue {
	cph := &ColorPickerHue{}
	cph.size = size

	cph.colors = [7]color.NRGBA{
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 255, 255, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{255, 0, 0, 255},
	}

	cph.partLength = float32(cph.size.X) / float32(len(cph.colors)-1)
	cph.updateChosenColorFromPickerPos(0.0)

	return cph
}

func (cph *ColorPickerHue) Layout(gtx layout.Context) layout.Dimensions {
	cph.triggerRenderImageUpdate = cph.triggerRenderImageUpdate || cph.size != gtx.Constraints.Max
	cph.size = gtx.Constraints.Max
	cph.partLength = float32(cph.size.X) / float32(len(cph.colors)-1)
	
	cph.HandleInput(gtx)
	cph.Draw(gtx)
	
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (cph *ColorPickerHue) HandleInput(gtx layout.Context) {
	cph.pickedNewColor = false

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
		cph.pickedNewColor = true
	}

	area.Pop()
}

func (cph *ColorPickerHue) updateChosenColorFromPickerPos(fraction float32) {
	cph.pickFractionPos = fraction
	cph.pickFractionPos = float32(math.Max(math.Min(float64(cph.pickFractionPos), 1), 0))

	color := cph.getColorFromPosition(cph.getPickerPositionClamped())
	cph.triggerRenderImageUpdate = cph.triggerRenderImageUpdate || color != cph.chosenColor
	cph.chosenColor = color
}

func (cph *ColorPickerHue) getPickerPositionClamped() float32 {
	pickerPosition := cph.pickFractionPos * float32(cph.size.X)
	pickerPosition = float32(math.Max(math.Min(float64(pickerPosition), float64(cph.size.X)), 0))
	return pickerPosition
}

func (cph *ColorPickerHue) Draw(gtx layout.Context) {
	grect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

	if cph.triggerRenderImageUpdate {
		cph.triggerRenderImageUpdate = false
		img := image.NewNRGBA(grect)
		for x := 0; x < grect.Dx(); x++ {
			col := cph.getColorFromPosition(float32(x))
			for y := 0; y < grect.Dy(); y++ {
				img.SetNRGBA(x, y, col)
			}
		}
		cph.renderImage = paint.NewImageOp(img)
	}

	cph.renderImage.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	p := cph.getPickerPositionClamped()
	col := cph.getColorFromPosition(p)
	drawPicker(f32.Pt(p, float32(cph.size.Y)/2), col, gtx)
}

func (cph *ColorPickerHue) getColorFromPosition(x float32) color.NRGBA {
	gradientNumber := utils.ClampInt(int(x/cph.partLength), 0, len(cph.colors)-2)
	col1 := cph.colors[gradientNumber]
	col2 := cph.colors[gradientNumber+1]
	gradientStride := float64(x-cph.partLength*float32(gradientNumber)) / float64(cph.partLength)
	colorResult := lerpColor(col1, col2, gradientStride)
	return colorResult
}

func (cph *ColorPickerHue) getPositionFractionFromColor(col color.NRGBA) float32 {
	h, _, _ := utils.RgbToHsv(col.R, col.G, col.B)
	return float32(h)
}

func (cph *ColorPickerHue) updateColor(col color.NRGBA) {
	cph.updateChosenColorFromPickerPos(cph.getPositionFractionFromColor(col))
}
