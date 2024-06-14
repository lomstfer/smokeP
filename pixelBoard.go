package main

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

const (
	zoomMultiplier     = 0.01
	defaultBoardWidth  = 32
	defaultBoardHeight = 32
)

type PixelBoard struct {
	pixelImg               *image.NRGBA
	pixelImgOp             paint.ImageOp
	distanceMoved          f32.Point
	position               f32.Point
	scale                  float32
	drawingColor           color.NRGBA
	previousDrawPixelCoord *image.Point
}

func newPixelBoard() *PixelBoard {
	pb := &PixelBoard{}

	pb.setToNewImage(image.NewNRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight)))
	pb.pixelImgOp = paint.NewImageOp(pb.pixelImg)
	pb.pixelImgOp.Filter = paint.FilterNearest

	return pb
}

func (pb *PixelBoard) setToNewImage(newImage *image.NRGBA) {
	pb.pixelImg = newImage
	pb.pixelImgOp = paint.NewImageOp(pb.pixelImg)
	pb.pixelImgOp.Filter = paint.FilterNearest
	pb.scale = 640.0 / float32(math.Sqrt(float64(newImage.Rect.Dx()*newImage.Rect.Dx())+float64(newImage.Rect.Dy()*newImage.Rect.Dy())))
	pb.distanceMoved = pb.Size().Div(-2)
}

func (pb *PixelBoard) Size() f32.Point {
	return f32.Pt(float32(pb.pixelImgOp.Size().X)*pb.scale, float32(pb.pixelImgOp.Size().Y)*pb.scale)
}

func (pb *PixelBoard) Update(editingAreaCenter f32.Point) {
	pb.position = pb.distanceMoved.Add(editingAreaCenter)
}

func (pb *PixelBoard) DrawSelf(ops *op.Ops) {
	pb.pixelImgOp.Add(ops)

	intPosition := f32.Pt(float32(int(pb.position.X)), float32(int(pb.position.Y)))
	tStack := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(pb.scale, pb.scale)).Offset(intPosition)).Push(ops)
	paint.PaintOp{}.Add(ops)
	tStack.Pop()
}

func (pb *PixelBoard) OnDraw(mousePos f32.Point) {
	// size := pb.Size()
	// onBoard := mousePos.X > pb.position.X &&
	// 	mousePos.X < pb.position.X+size.X &&
	// 	mousePos.Y > pb.position.Y &&
	// 	mousePos.Y < pb.position.Y+size.Y

	rel := mousePos.Sub(pb.position).Div(pb.scale)
	pixelCoord := image.Pt(int(rel.X), int(rel.Y))

	if pb.previousDrawPixelCoord != nil {
		drawLineBetweenPoints(pb.pixelImg, *pb.previousDrawPixelCoord, pixelCoord, pb.drawingColor)
	} else {
		pb.pixelImg.SetNRGBA(pixelCoord.X, pixelCoord.Y, pb.drawingColor)
	}

	pb.pixelImgOp = paint.NewImageOp(pb.pixelImg)
	pb.pixelImgOp.Filter = paint.FilterNearest

	pb.previousDrawPixelCoord = &pixelCoord
}

func (pb *PixelBoard) OnStopDrawing() {
	pb.previousDrawPixelCoord = nil
}

func (pb *PixelBoard) Zoom(editingAreaCenter f32.Point, scrollY float32, mousePos f32.Point) {
	size := pb.Size()
	scaleChange := -scrollY * zoomMultiplier * pb.scale
	pb.scale += scaleChange

	mouseRelBoard := mousePos.Sub(pb.position)

	ratioX := mouseRelBoard.X / size.X
	ratioY := mouseRelBoard.Y / size.Y
	pb.distanceMoved = pb.distanceMoved.Sub(f32.Pt(
		ratioX*scaleChange*float32(pb.pixelImgOp.Size().X),
		ratioY*scaleChange*float32(pb.pixelImgOp.Size().Y),
	))

	pb.position = pb.distanceMoved.Add(editingAreaCenter)
}

func drawLineBetweenPoints(img *image.NRGBA, p1 image.Point, p2 image.Point, color color.NRGBA) {
	d := p2.Sub(p1)
	if int(math.Abs(float64(d.X))) > int(math.Abs(float64(d.Y))) {
		k := float32(d.Y) / float32(d.X)
		for i := 0; i <= int(math.Abs(float64(d.X))); i++ {
			x := i * int(float32(d.X)/float32(math.Abs(float64(d.X))))
			y := int(math.Round(float64(float32(x) * k)))
			pos := p1.Add(image.Pt(x, y))
			img.SetNRGBA(pos.X, pos.Y, color)
		}
	} else {
		k := float32(d.X) / float32(d.Y)
		for i := 0; i <= int(math.Abs(float64(d.Y))); i++ {
			y := i * int(float32(d.Y)/float32(math.Abs(float64(d.Y))))
			x := int(math.Round(float64(float32(y) * k)))
			pos := p1.Add(image.Pt(x, y))
			img.SetNRGBA(pos.X, pos.Y, color)
		}
	}
}
