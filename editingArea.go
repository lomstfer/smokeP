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

type EditingArea struct {
	mousePos f32.Point
	board    *PixelBoard
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()

	return ea
}

func (ea *EditingArea) Layout(gtx layout.Context) layout.Dimensions {
	// r := image.Rect(gtx.Constraints.Min.X, gtx.Constraints.Min.Y, gtx.Constraints.Max.X-gtx.Constraints.Min.X, gtx.Constraints.Max.Y-gtx.Constraints.Min.Y)
	// area := clip.Rect(r).Push(gtx.Ops)

	dragAccumulation := f32.Point{X: 0, Y: 0}

	event.Op(gtx.Ops, ea)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: ea,
			Kinds:  pointer.Drag | pointer.Move | pointer.Press | pointer.Scroll,
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Scroll:
			fmt.Println(e.Scroll)
		// case pointer.Drag:
		// 	if e.Buttons.Contain(pointer.ButtonSecondary) {
		// 		dragAccumulation = dragAccumulation.Add(e.Position.Sub(ea.mousePos))
		// 	}
		// 	if e.Buttons.Contain(pointer.ButtonPrimary) {
		// 		ea.board.CheckAndDraw(e.Position)
		// 	}
		// case pointer.Press:
		// 	if e.Buttons.Contain(pointer.ButtonPrimary) {
		// 		ea.board.CheckAndDraw(e.Position)
		// 	}
		}

		// ea.mousePos = e.Position
	}

	ea.board.position = ea.board.position.Add(dragAccumulation)

	// area.Pop()

	col := color.NRGBA{R: 255, A: 255}

	rectc := clip.Rect(image.Rect(int(ea.mousePos.X)-100, int(ea.mousePos.Y)-100, int(ea.mousePos.X)+100, int(ea.mousePos.Y)+100)).Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rectc.Pop()

	ea.board.Layout(gtx)

	return layout.Dimensions{Size: image.Pt(0, 0)}
}
