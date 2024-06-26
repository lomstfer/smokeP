package main

import (
	"image"
	"image/color"
	"smokep/pixeltools"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type EditingArea struct {
	mousePos f32.Point
	board    *PixelBoard
	size     image.Point
	center   f32.Point

	currentTool pixeltools.Tool
	pen         *pixeltools.Pencil
	picker      *bool
	pickedColor *color.NRGBA

	justChoseTool bool
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()
	ea.pen = &pixeltools.Pencil{}
	ea.picker = new(bool)
	ea.currentTool = pixeltools.ToolPen
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
	ea.UpdateTools(gtx)

	dragAccumulation := f32.Point{X: 0, Y: 0}
	ea.pickedColor = nil

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  ea,
			Kinds:   pointer.Drag | pointer.Press | pointer.Scroll | pointer.Leave | pointer.Release,
			ScrollY: pointer.ScrollRange{Min: -10, Max: 10},
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
				point := ea.board.GetPixelPoint(e.Position)

				switch ea.currentTool {
				case pixeltools.ToolPen:
					ea.pen.OnDraw(ea.board.pixelImg, point, ea.board.drawingColor)
					ea.board.refreshImage()
				}
			}
		case pointer.Press:
			if e.Buttons.Contain(pointer.ButtonPrimary) && !ea.justChoseTool {
				gtx.Execute(key.FocusCmd{Tag: ea})
				point := ea.board.GetPixelPoint(e.Position)

				switch ea.currentTool {
				case pixeltools.ToolPen:
					ea.pen.OnDraw(ea.board.pixelImg, point, ea.board.drawingColor)
					ea.board.refreshImage()
				case pixeltools.ToolPick:
					if ea.board.IsPointOnBoard(e.Position) {
						color := ea.board.pixelImg.NRGBAAt(point.X, point.Y)
						ea.pickedColor = &color
					}
				}
			}
			if e.Buttons.Contain(pointer.ButtonSecondary) {
				ea.mousePos = e.Position
			}
		case pointer.Leave, pointer.Release:
			switch ea.currentTool {
			case pixeltools.ToolPen:
				if action := ea.pen.OnEnd(); action != nil {
					ea.board.AddAction(action)
					ea.board.latestActionIndex += 1
				}
			}
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

func (ea *EditingArea) UpdateTools(gtx layout.Context) {
	ea.justChoseTool = false
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: ea.pen,
			Kinds:  pointer.Press,
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
		ea.currentTool = pixeltools.ToolPen
		ea.justChoseTool = true
	}
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: ea.picker,
			Kinds:  pointer.Press,
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
		ea.currentTool = pixeltools.ToolPick
		ea.justChoseTool = true
	}
}

func (ea *EditingArea) Layout(gtx layout.Context, gridBg *utils.GridBackground) layout.Dimensions {
	ea.size = gtx.Constraints.Max
	ea.center = f32.Pt(float32(ea.size.X), float32(ea.size.Y)).Div(2)
	defer clip.Rect(image.Rect(0, 0, ea.size.X, ea.size.Y)).Push(gtx.Ops).Pop()

	ea.Update(gtx)

	{
		boardPos := image.Pt(int(ea.board.position.X), int(ea.board.position.Y))
		boardSize := image.Pt(int(ea.board.Size().X), int(ea.board.Size().Y))
		r := image.Rect(boardPos.X, boardPos.Y, boardPos.X+boardSize.X, boardPos.Y+boardSize.Y)
		gridBg.Draw(gtx.Ops, r)
	}
	ea.board.DrawSelf(gtx.Ops)

	event.Op(gtx.Ops, ea)

	ea.LayoutTools(gtx)

	return layout.Dimensions{Size: ea.size}
}

func (ea *EditingArea) LayoutTools(gtx layout.Context) {
	size := 50
	space := 10
	increment := space

	// pen
	{
		r := image.Rect(space, increment, space+size, increment+size)
		a := clip.Rect(r).Push(gtx.Ops)
		DrawToolRect(gtx, r, ea.currentTool == pixeltools.ToolPen)
		event.Op(gtx.Ops, ea.pen)
		a.Pop()
	}
	increment += space + size

	// picker
	{
		r := image.Rect(space, increment, space+size, increment+size)
		a := clip.Rect(r).Push(gtx.Ops)
		DrawToolRect(gtx, r, ea.currentTool == pixeltools.ToolPick)
		event.Op(gtx.Ops, ea.picker)
		a.Pop()
	}
	increment += space + size
}

func DrawToolRect(gtx layout.Context, rect image.Rectangle, selected bool) {
	paint.ColorOp{Color: color.NRGBA{150, 150, 150, 255}}.Add(gtx.Ops)
	var s clip.Stack
	if selected {
		paint.PaintOp{}.Add(gtx.Ops)
		s = clip.Stroke{
			Path:  clip.RRect{Rect: rect}.Path(gtx.Ops),
			Width: 6,
		}.Op().Push(gtx.Ops)

		paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		s.Pop()
	} else {
		paint.PaintOp{}.Add(gtx.Ops)
	}
}
