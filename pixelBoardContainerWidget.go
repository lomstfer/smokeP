package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type PixelBoardContainerWidget struct {
	widget.BaseWidget
	boardHolder       *fyne.Container
	board             *PixelBoard
	wantsToPaint      bool
	dragMouseStartPos *fyne.Position
}

func newPixelBoardContainerWidget(boardHolder *fyne.Container, board *PixelBoard) *PixelBoardContainerWidget {
	b := &PixelBoardContainerWidget{boardHolder: boardHolder, board: board}
	b.ExtendBaseWidget(b)
	return b
}

func (pbc *PixelBoardContainerWidget) Scrolled(e *fyne.ScrollEvent) {
	pbc.board.Zoom(e.Scrolled.DY, pbc.positionRelBoard(e.Position))
}

func (pbc *PixelBoardContainerWidget) CreateRenderer() fyne.WidgetRenderer {
	// return &emptyRenderer{}
	return widget.NewSimpleRenderer(pbc.boardHolder)
}

func (pbc *PixelBoardContainerWidget) MouseDown(e *desktop.MouseEvent) {
	if e.Button == desktop.MouseButtonPrimary {
		pbc.wantsToPaint = true
		pbc.board.Paint(pbc.positionRelBoard(e.Position))
	}
	if e.Button == desktop.MouseButtonSecondary {
		p := fyne.NewPos(e.AbsolutePosition.X, e.AbsolutePosition.Y)
		pbc.dragMouseStartPos = &p
	}
}

func (pbc *PixelBoardContainerWidget) MouseUp(e *desktop.MouseEvent) {
	if e.Button == desktop.MouseButtonPrimary {
		pbc.wantsToPaint = false
	}
	if e.Button == desktop.MouseButtonSecondary {
		pbc.dragMouseStartPos = nil
		pbc.board.endMove()
	}
}

func (pbc *PixelBoardContainerWidget) MouseMoved(e *desktop.MouseEvent) {
	if pbc.wantsToPaint {
		pbc.board.Paint(pbc.positionRelBoard(e.Position))
	}
	if e.Button == desktop.MouseButtonSecondary && pbc.dragMouseStartPos != nil {
		pbc.board.updateDrag(e.AbsolutePosition, *pbc.dragMouseStartPos)
	}
}

func (pbc *PixelBoardContainerWidget) MouseIn(*desktop.MouseEvent) {}

func (pbc *PixelBoardContainerWidget) MouseOut() {
	pbc.wantsToPaint = false
	pbc.dragMouseStartPos = nil
	pbc.board.endMove()
}

func (pbc *PixelBoardContainerWidget) positionRelBoard(position fyne.Position) fyne.Position {
	boardPosition := pbc.board.boardObj.Position().Add(fyne.NewSize(pbc.Size().Width/2, pbc.Size().Height/2)).Subtract(fyne.NewSize(defaultBoardPixelWidth/2, defaultBoardPixelHeight/2))
	positionRelBoard := position.Subtract(boardPosition)
	return positionRelBoard
}

func (pbc *PixelBoardContainerWidget) isPositionOnBoard(position fyne.Position) bool {
	positionRelBoard := pbc.positionRelBoard(position)

	return positionRelBoard.X > 0 &&
		positionRelBoard.X < pbc.board.boardObj.Size().Width &&
		positionRelBoard.Y > 0 &&
		positionRelBoard.Y < pbc.board.boardObj.Size().Height
}
