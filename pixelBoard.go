package main

import (
	"image"
	"image/color"
	"math"
	"smokep/boardactions"
	"smokep/utils"

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
	previousDrawPixelPoint *image.Point
	actionList             []boardactions.Action
	latestActionIndex      int
	currentDrawAction      *boardactions.DrawAction
}

func newPixelBoard() *PixelBoard {
	pb := &PixelBoard{}

	pb.setToNewImage(image.NewNRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight)))
	pb.refreshImage()

	pb.latestActionIndex = -1

	return pb
}

func (pb *PixelBoard) refreshImage() {
	pb.pixelImgOp = paint.NewImageOp(pb.pixelImg)
	pb.pixelImgOp.Filter = paint.FilterNearest
}

func (pb *PixelBoard) setToNewImage(newImage *image.NRGBA) {
	pb.pixelImg = newImage
	pb.refreshImage()
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


func (pb *PixelBoard) AddAction(action boardactions.Action) {
	for i := len(pb.actionList) - 1; i > pb.latestActionIndex; i-- {
		pb.actionList = append(pb.actionList[:i], pb.actionList[i+1:]...)
	}

	pb.actionList = append(pb.actionList, action)
}

func (pb *PixelBoard) Undo() {
	if pb.latestActionIndex <= -1 {
		return
	}

	pb.actionList[pb.latestActionIndex].Undo(pb.pixelImg)
	pb.setToNewImage(pb.pixelImg)
	pb.latestActionIndex -= 1
}

func (pb *PixelBoard) Redo() {
	if len(pb.actionList) == 0 || pb.latestActionIndex == len(pb.actionList)-1 {
		return
	}

	pb.latestActionIndex += 1
	pb.actionList[pb.latestActionIndex].Do(pb.pixelImg)
	pb.setToNewImage(pb.pixelImg)
}

func (pb *PixelBoard) Resize(newSize image.Point, resizeOrigin f32.Point) {
	pb.AddAction(boardactions.NewResizeAction(newSize, resizeOrigin))
	pb.Redo()
}

func (pb *PixelBoard) OnDraw(mousePos f32.Point) {
	// size := pb.Size()
	// onBoard := mousePos.X > pb.position.X &&
	// 	mousePos.X < pb.position.X+size.X &&
	// 	mousePos.Y > pb.position.Y &&
	// 	mousePos.Y < pb.position.Y+size.Y

	rel := mousePos.Sub(pb.position).Div(pb.scale)
	pixelPoint := image.Pt(int(rel.X), int(rel.Y))

	if pb.pixelImg.NRGBAAt(pixelPoint.X, pixelPoint.Y) == pb.drawingColor {
		return
	}

	if pb.currentDrawAction == nil {
		pb.currentDrawAction = boardactions.NewDrawAction(nil, pb.drawingColor)
	}

	if pb.previousDrawPixelPoint != nil {
		var points []image.Point
		{
			betweenPoints := utils.GetLineBetweenPoints(*pb.previousDrawPixelPoint, pixelPoint)
			for _, p := range betweenPoints {
				if pb.pixelImg.NRGBAAt(p.X, p.Y) != pb.drawingColor {
					points = append(points, p)
				}
			}
		}
		
		for _, p := range points {
			pb.currentDrawAction.PreviousPixelcolors[p] = pb.pixelImg.NRGBAAt(p.X, p.Y)
			pb.pixelImg.SetNRGBA(p.X, p.Y, pb.drawingColor)
		}
		pb.currentDrawAction.PixelPoints = append(pb.currentDrawAction.PixelPoints, points...)
	} else {
		pb.currentDrawAction.PreviousPixelcolors[pixelPoint] = pb.pixelImg.NRGBAAt(pixelPoint.X, pixelPoint.Y)
		pb.pixelImg.SetNRGBA(pixelPoint.X, pixelPoint.Y, pb.drawingColor)
		pb.currentDrawAction.PixelPoints = append(pb.currentDrawAction.PixelPoints, pixelPoint)
	}

	pb.refreshImage()

	pb.previousDrawPixelPoint = &pixelPoint
}

func (pb *PixelBoard) OnStopDrawing() {
	pb.previousDrawPixelPoint = nil
	if pb.currentDrawAction != nil {
		pb.AddAction(*pb.currentDrawAction)
		pb.latestActionIndex += 1
	}
	pb.currentDrawAction = nil
}
