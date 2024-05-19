package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
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
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
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

	img := image.NewRGBA(image.Rect(0, 0, defaultBoardPixelWidth, defaultBoardPixelHeight))

	for i := range img.Pix {
		img.Pix[i] = uint8(rand.Intn(255))
	}

	imgObj := canvas.NewImageFromImage(img)
	imgObj.ScaleMode = canvas.ImageScalePixels

	pixelBoard := newPixelBoard(imgObj)
	pixelBoard.img.FillMode = canvas.ImageFillOriginal
	// pixelBoard.Zoom(0, fyne.NewPos(0, 0))
	
	pixelBoard.img.Resize(fyne.NewSize(defaultBoardWidth, defaultBoardHeight))

	pixelBoardHolder := container.NewWithoutLayout(pixelBoard.img)
	pixelBoardHolderCenterer := container.NewCenter(pixelBoardHolder)
	pixelBoardContainerWidget := newPixelBoardContainerWidget(pixelBoardHolderCenterer, pixelBoard)

	fmt.Println("e")

	leftPanel := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
	rightPanel := canvas.NewRectangle(color.RGBA{0, 0, 255, 255})

	cont := container.NewGridWithColumns(3)
	cont.Add(leftPanel)
	cont.Add(pixelBoardContainerWidget)
	cont.Add(rightPanel)

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
		}
	}()

	window.SetContent(cont)
	window.ShowAndRun()
}
