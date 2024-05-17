package main

import (
	"image"
	"image/color"
	"time"

	"math/rand"


	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"

	"fyne.io/fyne/v2/widget"
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

type clickableImage struct {
	widget.BaseWidget
	img *canvas.Image
	pressed bool
}

func newClickableImage(img *canvas.Image) *clickableImage {
	c := &clickableImage{img: img}
	c.ExtendBaseWidget(c)
	return c
}

func (d *clickableImage) MouseDown(e *desktop.MouseEvent) {
	d.pressed = true

	pixelPosX := int(e.Position.X / d.Size().Width * float32(d.img.Image.Bounds().Dx()))
	pixelPosY := int(e.Position.Y / d.Size().Height * float32(d.img.Image.Bounds().Dy()))

	rgbaImg := getPixelData(d.img.Image)
	rgbaImg.Set(pixelPosX, pixelPosY, color.RGBA{255, 255, 255, 255})
	d.img.Image = rgbaImg
	d.img.Refresh()

}

func (d *clickableImage) MouseUp(e *desktop.MouseEvent) {
	d.pressed = false
}

func (d *clickableImage) MouseMoved(e *desktop.MouseEvent) {
	if (!d.pressed) {
		return
	}

	pixelPosX := int(e.Position.X / d.Size().Width * float32(d.img.Image.Bounds().Dx()))
	pixelPosY := int(e.Position.Y / d.Size().Height * float32(d.img.Image.Bounds().Dy()))

	rgbaImg := getPixelData(d.img.Image)
	rgbaImg.Set(pixelPosX, pixelPosY, color.RGBA{255, 255, 255, 255})
	d.img.Image = rgbaImg
	d.img.Refresh()
}

func (d *clickableImage) MouseIn(*desktop.MouseEvent) {}

func (d *clickableImage) MouseOut() {
	d.pressed = false
}

func (c *clickableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.img)
}

// func (c *clickableImage) MinSize() fyne.Size {
//     // return c.img.Size()
// 	fmt.Println(c.img.Size())
// 	return fyne.NewSize(100, 100)
// }

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Pixel Click")
	window.Resize(fyne.NewSize(640, 360))

	const (
		defaultBoardWidth  = 8
		defaultBoardHeight = 8
	)

	img := image.NewRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight))

	for i, _ := range img.Pix {
		img.Pix[i] = uint8(rand.Intn(255))
	}

	imgObj := canvas.NewImageFromImage(img)
	imgObj.ScaleMode = canvas.ImageScalePixels

	clickableImg := newClickableImage(imgObj)
	clickableImg.img.FillMode = canvas.ImageFillOriginal

	clickableImg.Resize(fyne.NewSize(300, 300))

	cont := container.NewWithoutLayout(clickableImg)

	go func() {
		widthBefore := window.Content().Size().Width
		heightBefore := window.Content().Size().Height
		for {
			width := window.Content().Size().Width
			height := window.Content().Size().Height
			if width != widthBefore || height != heightBefore {
				clickableImg.Move(fyne.NewPos(
					width/2-clickableImg.Size().Width/2,
					height/2-clickableImg.Size().Height/2,
				))
			}
			widthBefore = width
			heightBefore = height
			time.Sleep(20 * time.Millisecond) // you may want to change this
		}
	}()

	window.SetContent(cont)
	window.ShowAndRun()
}

// func main() {
// 	myApp := app.New()
// 	w := myApp.NewWindow("o")

// 	const (
// 		defaultBoardWidth  = 32
// 		defaultBoardHeight = 32
// 	)

// 	img := image.NewRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight))

// 	for i, _ := range img.Pix {
// 		img.Pix[i] = uint8(rand.Intn(255))
// 	}

// 	img.Set(50, 50, color.RGBA{255, 0, 0, 255})
// 	img.Set(60, 60, color.RGBA{0, 255, 0, 255})

// 	fyneImage := canvas.NewImageFromImage(img)
// 	fyneImage.FillMode = canvas.ImageFillOriginal
// 	fyneImage.ScaleMode = canvas.ImageScalePixels

// 	w.SetContent(fyneImage)

// 	w.Resize(fyne.NewSize(120, 100))
// 	w.ShowAndRun()
// }
