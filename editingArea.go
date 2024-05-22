package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

type EditingArea struct {
	mousePos      f32.Point
	board         *PixelBoard
	size          image.Point
	boardCentered bool
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()

	return ea
}

func (ea *EditingArea) CenterBoard() {
	rounded := ea.size.Div(2).Sub(ea.board.Size().Div(2).Round())
	ea.board.position = f32.Pt(float32(rounded.X), float32(rounded.Y))
}

func (ea *EditingArea) Layout(gtx layout.Context) layout.Dimensions {
	ea.size = gtx.Constraints.Max
	if !ea.boardCentered {
		ea.CenterBoard()
		ea.boardCentered = true
	}
	area := clip.Rect(image.Rect(0, 0, ea.size.X, ea.size.Y)).Push(gtx.Ops)

	dragAccumulation := f32.Point{X: 0, Y: 0}

	event.Op(gtx.Ops, ea)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       ea,
			Kinds:        pointer.Drag | pointer.Move | pointer.Press | pointer.Scroll,
			ScrollBounds: image.Rect(-10, -10, 10, 10),
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
			ea.board.Zoom(e.Scroll.Y, e.Position)
		case pointer.Drag:
			if e.Buttons.Contain(pointer.ButtonSecondary) {
				dragAccumulation = dragAccumulation.Add(e.Position.Sub(ea.mousePos))
			}
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				ea.board.CheckIfOnBoardAndDraw(e.Position)
			}
		case pointer.Press:
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				ea.board.CheckIfOnBoardAndDraw(e.Position)
			}
		}

		ea.mousePos = e.Position
	}

	ea.board.position = ea.board.position.Add(dragAccumulation)

	ea.board.Draw(gtx.Ops)

	area.Pop()

	return layout.Dimensions{Size: image.Pt(ea.size.X, ea.size.Y)}
}
