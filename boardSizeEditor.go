package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type BoardSizeEditor struct {
	editor         *widget.Editor
	submitResult   *image.Point
	selectedOrigin f32.Point
	squares        [3][3]*bool
}

func NewBoardSizeEditor(pixelBoardSize image.Point) *BoardSizeEditor {
	bse := &BoardSizeEditor{}
	bse.editor = &widget.Editor{Submit: true}
	bse.editor.SetText(fmt.Sprintf("%dx%d", pixelBoardSize.X, pixelBoardSize.Y))
	bse.squares = [3][3]*bool{
		{new(bool), new(bool), new(bool)},
		{new(bool), new(bool), new(bool)},
		{new(bool), new(bool), new(bool)},
	}

	return bse
}

func (bse *BoardSizeEditor) Update(gtx layout.Context, pixelBoardSize image.Point) {
	if !gtx.Focused(bse.editor) {
		bse.editor.SetText(fmt.Sprintf("%dx%d", pixelBoardSize.X, pixelBoardSize.Y))
	}

	for y := 0; y < len(bse.squares); y++ {
		for x := 0; x < len(bse.squares[y]); x++ {
			for {
				ev, ok := gtx.Event(pointer.Filter{
					Target: bse.squares[y][x],
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

				bse.selectedOrigin.X = float32(x) / 2
				bse.selectedOrigin.Y = float32(y) / 2
				gtx.Execute(key.FocusCmd{Tag: bse.squares[y][x]})
			}
		}
	}

	bse.submitResult = nil
	for {
		ev, ok := bse.editor.Update(gtx)
		if !ok {
			break
		}
		_, ok = ev.(widget.SubmitEvent)
		if !ok {
			continue
		}
		gtx.Execute(key.FocusCmd{Tag: nil})

		{
			input := bse.editor.Text()
			var width, height int
			_, err := fmt.Sscanf(input, "%dx%d", &width, &height)
			if err == nil {
				pt := image.Pt(width, height)
				bse.submitResult = &pt
			}
		}
	}
}

func (bse *BoardSizeEditor) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	d := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			d := material.Editor(theme, bse.editor, "").Layout(gtx)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			totalSize := 81
			iSize := totalSize / 3
			for y := 0; y < 3; y++ {
				for x := 0; x < 3; x++ {
					rect := image.Rect(x*iSize, y*iSize, x*iSize+iSize, y*iSize+iSize)
					area := clip.Rect(rect).Push(gtx.Ops)
					if float32(x)/2.0 == float32(bse.selectedOrigin.X) && float32(y)/2.0 == float32(bse.selectedOrigin.Y) {
						paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
					}
					s := clip.Stroke{
						Path:  clip.RRect{Rect: rect}.Path(gtx.Ops),
						Width: 3,
					}.Op().Push(gtx.Ops)
					paint.ColorOp{Color: color.NRGBA{255, 255, 255, 255}}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					s.Pop()
					event.Op(gtx.Ops, bse.squares[y][x])
					area.Pop()
				}
			}
			dims := layout.Dimensions{Size: image.Pt(totalSize, totalSize)}
			return dims
		}),
	)
	return d
}
