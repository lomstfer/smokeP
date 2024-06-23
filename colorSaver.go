package main

import (
	"image"
	"image/color"
	"smokep/utils"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ColorSaver struct {
	colors           []*color.NRGBA
	ws               *widget.Scrollbar
	saveColorButton  *widget.Clickable
	setColor         chan color.NRGBA
	justClickedColor *color.NRGBA
}

func NewColorSaver() *ColorSaver {
	cs := &ColorSaver{}
	cs.ws = &widget.Scrollbar{}
	cs.saveColorButton = &widget.Clickable{}
	cs.setColor = make(chan color.NRGBA)
	return cs
}

func (cs *ColorSaver) Update(gtx layout.Context, currentColor color.NRGBA) {
	if cs.saveColorButton.Clicked(gtx) {
		cs.colors = append(cs.colors, &currentColor)
	}

	cs.justClickedColor = nil
	for i := len(cs.colors) - 1; i >= 0; i-- {
		color := cs.colors[i]
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: color,
				Kinds:  pointer.Press,
			})
			if !ok {
				break
			}

			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			if e.Buttons.Contain(pointer.ButtonPrimary) {
				cs.justClickedColor = color
			}
			if e.Buttons.Contain(pointer.ButtonSecondary) {
				cs.colors = append(cs.colors[:i], cs.colors[i+1:]...)
			}
		}
	}
}

func (cs *ColorSaver) Layout(gtx layout.Context, theme *material.Theme, gridBg *utils.GridBackground) layout.Dimensions {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			d := material.Button(theme, cs.saveColorButton, "Save current color").Layout(gtx)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			r := image.Rect(0, 0, gtx.Constraints.Max.X, 70)
			area52 := clip.Rect(r).Push(gtx.Ops)
			posx := 0
			for i := len(cs.colors) - 1; i >= 0; i-- {
				c := cs.colors[i]
				rect := image.Rect(posx, 10, posx+50, 60)
				gridBg.Draw(gtx.Ops, rect)
				area := clip.Rect(rect).Push(gtx.Ops)
				event.Op(gtx.Ops, c)
				paint.ColorOp{Color: *c}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				area.Pop()
				posx += 60
			}
			area52.Pop()
			return layout.Dimensions{Size: image.Pt(r.Dx(), r.Dy())}
		}),
	)
	return d
}
