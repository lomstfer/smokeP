package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type EditingArea struct {
	mousePos        f32.Point
	board           *PixelBoard
	size            image.Point
	center          f32.Point
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()
	return ea
}

func (ea *EditingArea) Layout(gtx layout.Context) layout.Dimensions {
	ea.size = gtx.Constraints.Max
	ea.center = f32.Pt(float32(ea.size.X), float32(ea.size.Y)).Div(2)
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
			ea.board.Zoom(ea.center, e.Scroll.Y, e.Position)
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

	ea.board.distanceMoved = ea.board.distanceMoved.Add(dragAccumulation)

	ea.board.Update(ea.center)

	ea.drawBoardBackground(gtx.Ops)

	ea.board.Draw(gtx.Ops)

	area.Pop()

	return layout.Dimensions{Size: image.Pt(ea.size.X, ea.size.Y)}
}

func (ea *EditingArea) drawBoardBackground(ops *op.Ops) {
	boardIntPosition := image.Pt(int(ea.board.position.X), int(ea.board.position.Y))
	bsize := image.Pt(int(ea.board.Size().X), int(ea.board.Size().Y))
	bgCol := g_theme.Bg

	{
		area := clip.Rect(image.Rect(0, 0, ea.size.X, boardIntPosition.Y)).Push(ops)
		paint.ColorOp{Color: bgCol}.Add(ops)
		paint.PaintOp{}.Add(ops)
		area.Pop()
	}
	{
		area := clip.Rect(image.Rect(0, boardIntPosition.Y + bsize.Y, ea.size.X, ea.size.Y)).Push(ops)
		paint.ColorOp{Color: bgCol}.Add(ops)
		paint.PaintOp{}.Add(ops)
		area.Pop()
	}


	{
		area := clip.Rect(image.Rect(0, boardIntPosition.Y, boardIntPosition.X, boardIntPosition.Y + bsize.Y)).Push(ops)
		paint.ColorOp{Color: bgCol}.Add(ops)
		paint.PaintOp{}.Add(ops)
		area.Pop()
	}
	{
		area := clip.Rect(image.Rect(boardIntPosition.X + bsize.X, boardIntPosition.Y, ea.size.X, boardIntPosition.Y + bsize.Y)).Push(ops)
		paint.ColorOp{Color: bgCol}.Add(ops)
		paint.PaintOp{}.Add(ops)
		area.Pop()
	}
}
