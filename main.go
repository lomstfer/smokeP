package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	defaultCanvasWidth  = 8
	defaultCanvasHeight = 8
)

var defaultDrawColor = color.RGBA{255, 255, 255, 255}

type Vec2i struct {
	x int
	y int
}

type App struct {
	screen     *ebiten.Image
	canvas     *ebiten.Image
	pixels     map[Vec2i]color.Color
	background *ebiten.Image
}

func screenToCanvasPosition(x int, y int, screenWidth int, screenHeight int, canvasWidth int, canvasHeight int) (int, int) {
	pixelWidth := float64(screenWidth) / float64(canvasWidth)
	pixelHeight := float64(screenHeight) / float64(canvasHeight)

	x = int(math.Round((float64(x) - pixelWidth/2) / float64(screenWidth) * float64(canvasWidth)))
	y = int(math.Round((float64(y) - pixelHeight/2) / float64(screenHeight) * float64(canvasHeight)))
	return x, y
}

// func canvasToScreenPosition(x int, y int, canvasWidth int, canvasHeight int, screenWidth int, screenHeight int) (int, int) {
// 	x = int(math.Round(float64(x) / float64(canvasWidth) * float64(screenWidth)))
// 	y = int(math.Round(float64(y) / float64(canvasHeight) * float64(screenHeight)))
// 	return x, y
// }

func getPixelFromCursorPosition(pixels *map[Vec2i]color.Color, screenWidth int, screenHeight int, canvasWidth int, canvasHeight int) color.Color {
	x, y := ebiten.CursorPosition()
	x, y = screenToCanvasPosition(x, y, screenWidth, screenHeight, canvasWidth, canvasHeight)

	color := (*pixels)[Vec2i{x, y}]
	return color
}

func setPixelAtPosition(x int, y int, newColor color.Color, pixels *map[Vec2i]color.Color, screenWidth int, screenHeight int, canvasWidth int, canvasHeight int) (success bool) {
	x, y = screenToCanvasPosition(x, y, screenWidth, screenHeight, canvasWidth, canvasHeight)
	if x < 0 || y < 0 || x > canvasWidth-1 || y > canvasHeight-1 {
		return false
	}

	(*pixels)[Vec2i{x, y}] = newColor
	return true
}

func setDefaultPixels(pixels *map[Vec2i]color.Color, canvasWidth int, canvasHeight int, defaultColor color.Color) {
	x, y := 0, -1
	for i := 0; i < canvasWidth*canvasHeight; i++ {
		x = i % canvasWidth
		if x == 0 {
			y += 1
		}
		(*pixels)[Vec2i{x, y}] = defaultColor
	}
}

func (app *App) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && app.screen != nil {
		x, y := ebiten.CursorPosition()
		setPixelAtPosition(x, y, defaultDrawColor, &app.pixels, app.screen.Bounds().Dx(), app.screen.Bounds().Dy(), app.canvas.Bounds().Dx(), app.canvas.Bounds().Dy())
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) && app.screen != nil {
		x, y := ebiten.CursorPosition()
		setPixelAtPosition(x, y, color.RGBA{255, 0, 0, 255}, &app.pixels, app.screen.Bounds().Dx(), app.screen.Bounds().Dy(), app.canvas.Bounds().Dx(), app.canvas.Bounds().Dy())
	}

	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	app.screen = screen

	{
		bgGeo := ebiten.GeoM{}
		bgGeo.Scale(float64(screen.Bounds().Dx()) / float64(app.background.Bounds().Dx()), float64(screen.Bounds().Dy()) / float64(app.background.Bounds().Dy()))
		screen.DrawImage(app.background, &ebiten.DrawImageOptions{GeoM: bgGeo})
	}

	cw := app.canvas.Bounds().Dx()
	ch := app.canvas.Bounds().Dy()

	{
		x, y := ebiten.CursorPosition()
		vector.DrawFilledCircle(screen, float32(x), float32(y), 10, color.RGBA{255, 0, 0, 255}, false)
	}

	for i := range app.pixels {
		vector.DrawFilledRect(app.canvas, float32(i.x), float32(i.y), 1, 1, app.pixels[i], false)
	}

	{
		canvasGeo := ebiten.GeoM{}
		canvasGeo.Scale(float64(screen.Bounds().Dx()) / float64(cw), float64(screen.Bounds().Dy()) / float64(ch))
		screen.DrawImage(app.canvas, &ebiten.DrawImageOptions{GeoM: canvasGeo})
	}

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (app *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	background, _, err := ebitenutil.NewImageFromFile("test.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	img := ebiten.NewImage(defaultCanvasWidth, defaultCanvasHeight)

	app := App{
		nil,
		img,
		make(map[Vec2i]color.Color, defaultCanvasWidth*defaultCanvasHeight),
		background,
	}


	setDefaultPixels(&app.pixels, defaultCanvasWidth, defaultCanvasHeight, color.RGBA{0, 0, 0, 0})


	if err := ebiten.RunGame(&app); err != nil {
		log.Fatal(err)
	}
}
