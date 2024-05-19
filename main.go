package main

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/lusingander/colorpicker"
)

func getPixelData(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())

	for y := 0; y < img.Bounds().Dx(); y++ {
		for x := 0; x < img.Bounds().Dy(); x++ {
			c := img.At(x, y)
			rgba.Set(x, y, c)
		}
	}

	return rgba
}

var g_keysPressed = make(map[fyne.KeyName]bool)

func main() {
	myApp := app.NewWithID("smokep")
	window := myApp.NewWindow("smokep")
	window.Resize(fyne.NewSize(640, 360))

	if deskCanvas, ok := window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
			g_keysPressed[key.Name] = true
		})
		deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
			delete(g_keysPressed, key.Name)
		})
	}
	
	pixelBoard := newPixelBoard(canvas.NewImageFromFile("test.jpeg"))
	pixelBoardHolder := container.NewWithoutLayout(pixelBoard.boardObj)
	pixelBoardHolderCenterer := container.NewCenter(pixelBoardHolder)
	pixelBoardPanel := newPixelBoardContainerWidget(pixelBoardHolderCenterer, pixelBoard)

	leftPanel := container.NewStack()
	leftPanel.Add(canvas.NewRectangle(color.RGBA{255,0,0,255}))
	
	picker := colorpicker.New(200 /* height */, colorpicker.StyleHue /* Style */)
	picker.SetOnChanged(func(c color.Color) {
		pixelBoard.paintingColor = c
	})

	leftPanel.Add(picker)

	rightPanel := canvas.NewRectangle(color.RGBA{0, 0, 255, 255})

	mainPanels := container.NewGridWithColumns(3)
	mainPanels.Add(leftPanel)
	mainPanels.Add(pixelBoardPanel)
	mainPanels.Add(rightPanel)

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
		}
	}()

	window.SetContent(mainPanels)
	window.ShowAndRun()
}
