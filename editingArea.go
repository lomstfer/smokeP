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
	boardBackground *image.NRGBA
}

func newEditingArea() *EditingArea {
	ea := &EditingArea{}
	ea.board = newPixelBoard()
	ea.boardBackground = generateBoardBackground()
	return ea
}

func generateBoardBackground() *image.NRGBA {
	const pixelRow = 100
	img := image.NewNRGBA(image.Rect(0, 0, pixelRow, pixelRow))
	light := false
	for i := 0; i < len(img.Pix); i += 4 {
		if (i / 4 % pixelRow == 0) {
			light = !light
		}

		if light {
			img.Pix[i] = 200
			img.Pix[i + 1] = 200
			img.Pix[i + 2] = 200
			img.Pix[i + 3] = 255
		} else {
			img.Pix[i] = 100
			img.Pix[i + 1] = 100
			img.Pix[i + 2] = 100
			img.Pix[i + 3] = 255
		}

		light = !light
	}

	return img
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

	ea.drawBoardBackground(gtx.Ops)

	ea.board.Draw(ea.center, gtx.Ops)

	area.Pop()

	return layout.Dimensions{Size: image.Pt(ea.size.X, ea.size.Y)}
}

func (ea *EditingArea) drawBoardBackground(ops *op.Ops) {
	bpos := ea.board.position.Round()
	bsize := ea.board.Size().Round()
	area := clip.Rect(image.Rect(bpos.X, bpos.Y, bpos.X + bsize.X, bpos.Y + bsize.Y)).Push(ops)

	boardBgImgOp := paint.NewImageOp(ea.boardBackground)
	boardBgImgOp.Filter = paint.FilterNearest
	boardBgImgOp.Add(ops)

	scale := max(float32(ea.size.X) / float32(boardBgImgOp.Size().X), float32(ea.size.Y) / float32(boardBgImgOp.Size().Y))
	pos := ea.center.Sub(f32.Pt(float32(boardBgImgOp.Size().X) * scale, float32(boardBgImgOp.Size().Y) * scale).Div(2))
	tStack := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale)).Offset(pos)).Push(ops)
	paint.PaintOp{}.Add(ops)
	tStack.Pop()

	area.Pop()
}