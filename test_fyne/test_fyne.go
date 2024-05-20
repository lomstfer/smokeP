package main

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.NewWithID("s")
	window := myApp.NewWindow("s")
	window.Resize(fyne.NewSize(640, 360))

	test := newtestWidg()

	c := container.NewWithoutLayout(test)
	c.Resize(fyne.NewSize(300, 300))


	go func() {
		for {
			time.Sleep(time.Duration(694) * time.Microsecond)
			t := 1 - math.Exp(-0.01 * 16.666)
			lerpx := lerp(test.Position().X, test.wants.X, float32(t))
			lerpy := lerp(test.Position().Y, test.wants.Y, float32(t))
			test.Move(fyne.NewPos(lerpx, lerpy))
			// Vec2.lerp(velocity, new Vec2(0, 0), 1 - Math.exp(-CONSTS.MOVE_DAMP_SPEED * deltaTime))
		}
	}()

	window.SetContent(c)
	window.ShowAndRun()
}

type testWidget struct {
	widget.BaseWidget
	rect canvas.Rectangle
	lastMove time.Time
	wants fyne.Position
}

func newtestWidg() *testWidget {
	b := &testWidget{}
	b.Resize(fyne.NewSize(200, 200))
	b.rect = *canvas.NewRectangle(color.RGBA{255,0,0,255})
	b.ExtendBaseWidget(b)
	return b
}

func (tw *testWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(&tw.rect)
}

func (tw *testWidget) MouseDown(e *desktop.MouseEvent) {}

func (tw *testWidget) MouseUp(e *desktop.MouseEvent) {}

func (tw *testWidget) MouseMoved(e *desktop.MouseEvent) {
	// fmt.Println(time.Since(tw.lastMove))
	tw.wants = e.AbsolutePosition.Subtract(fyne.NewPos(100, 100))
	// tw.Move(e.AbsolutePosition.Subtract(fyne.NewPos(100, 100)))
	tw.Refresh()
	tw.lastMove = time.Now()
}

func (tw *testWidget) MouseIn(*desktop.MouseEvent) {}

func (tw *testWidget) MouseOut() {}

func lerp(p0 float32, p1 float32, t float32) float32 {
	if (t < 0) {
		t = 0
	}
	if (t > 1) {
		t = 1
	}
    return p0 + (p1 - p0) * t
}