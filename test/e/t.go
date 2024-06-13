package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Size(1280, 720), app.Title(""))

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	var ops op.Ops

	colorSelector := ColorSelector{size: image.Pt(0, 0)}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			
			colorSelector.Update(gtx) // removing this makes the color updates be one click behind

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					paint.ColorOp{Color: colorSelector.chosenColor}.Add(gtx.Ops)
					area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Push(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					area.Pop()
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return colorSelector.Layout(gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

type ColorSelector struct {
	size image.Point
	chosenColor color.NRGBA
}

func (cs *ColorSelector) Update(gtx layout.Context) {
	r := image.Rect(0, 0, cs.size.X, cs.size.Y)
	area := clip.Rect(r).Push(gtx.Ops)
	event.Op(gtx.Ops, cs)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       cs,
			Kinds:        pointer.Press,
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

		cs.chosenColor = color.NRGBA{0, uint8(e.Position.X / float32(cs.size.X) * 255), 0, 255}
	}

	area.Pop()
}

func (cs *ColorSelector) Layout(gtx layout.Context) layout.Dimensions {
	cs.size = gtx.Constraints.Max
	cs.Update(gtx) // This creates another input area, that is in the correct position. Without it the size is wrong (0, 0) at startup
	paint.ColorOp{Color: color.NRGBA{255, 0, 0, 255}}.Add(gtx.Ops)
	area := clip.Rect(image.Rect(0, 0, cs.size.X, cs.size.Y)).Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	area.Pop()
	return layout.Dimensions{Size: cs.size}
}
