package main

import (
	"image"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

type EditingArea struct {
	mousePos f32.Point
	board    *PixelBoard
	size     image.Point
	center   f32.Point
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()
	return ea
}

func (ea *EditingArea) Update(gtx layout.Context) {
	for {
		ev, ok := gtx.Event(key.Filter{
			Focus:    nil,
			Required: key.ModCtrl,
			Name:     "R",
		})
		if !ok {
			break
		}

		e, ok := ev.(key.Event)
		if !ok {
			continue
		}

		if e.State == key.Release {
			continue
		}

		ea.board.Resize(ea.board.pixelImgOp.Size().Add(image.Pt(10, 10)), f32.Pt(0.5, 0.5))
	}

	ea.CheckUndoRedo(gtx)

	dragAccumulation := f32.Point{X: 0, Y: 0}

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       ea,
			Kinds:        pointer.Drag | pointer.Press | pointer.Scroll | pointer.Leave | pointer.Release,
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
				ea.mousePos = e.Position
			}
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				ea.board.OnDraw(e.Position)
			}
		case pointer.Press:
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				gtx.Execute(key.FocusCmd{Tag: ea})
				ea.board.OnDraw(e.Position)
			}
			if e.Buttons.Contain(pointer.ButtonSecondary) {
				ea.mousePos = e.Position
			}
		case pointer.Leave, pointer.Release:
			ea.board.OnStopDrawing()
		}
	}

	ea.board.distanceMoved = ea.board.distanceMoved.Add(dragAccumulation)

	ea.board.Update(ea.center)
}

func (ea *EditingArea) CheckUndoRedo(gtx layout.Context) {
	for {
		ev, ok := gtx.Event(key.Filter{
			Focus:    nil,
			Required: key.ModCtrl,
			Optional: key.ModShift,
			Name:     "Z",
		})
		if !ok {
			break
		}

		e, ok := ev.(key.Event)
		if !ok {
			continue
		}

		if e.State == key.Release {
			continue
		}

		if e.Modifiers.Contain(key.ModShift) {
			ea.board.Redo()
		} else {
			ea.board.Undo()
		}
	}

	for {
		ev, ok := gtx.Event(key.Filter{
			Focus:    nil,
			Required: key.ModCtrl,
			Name:     "Y",
		})
		if !ok {
			break
		}

		e, ok := ev.(key.Event)
		if !ok {
			continue
		}

		if e.State == key.Release {
			continue
		}

		ea.board.Redo()
	}
}

func (ea *EditingArea) Layout(gtx layout.Context, gridBg *utils.GridBackground) layout.Dimensions {
	{
		area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Push(gtx.Ops)
		event.Op(gtx.Ops, ea)
		area.Pop()
	}

	ea.size = gtx.Constraints.Max
	ea.center = f32.Pt(float32(ea.size.X), float32(ea.size.Y)).Div(2)
	defer clip.Rect(image.Rect(0, 0, ea.size.X, ea.size.Y)).Push(gtx.Ops).Pop()

	ea.Update(gtx)

	{
		boardPos := ea.board.position.Round()
		boardSize := ea.board.Size().Round()
		r := image.Rect(boardPos.X, boardPos.Y, boardPos.X + boardSize.X, boardPos.Y + boardSize.Y)
		gridBg.Draw(gtx.Ops, r)
	}
	ea.board.DrawSelf(gtx.Ops)

	event.Op(gtx.Ops, ea)

	return layout.Dimensions{Size: ea.size}
}
