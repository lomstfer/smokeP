package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib moving square")
	defer rl.CloseWindow()

	rl.SetTargetFPS(144)

	// lastTime := time.Now()

	for !rl.WindowShouldClose() {
		// fmt.Println(time.Since(lastTime))
		// lastTime = time.Now()
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		mouse := rl.GetMousePosition()

		rl.DrawRectangle(int32(mouse.X) - 100, int32(mouse.Y) - 100, 200, 200, color.RGBA{255,0,0,255})

		rl.EndDrawing()
	}
}
