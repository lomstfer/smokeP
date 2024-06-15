package boardactions

import "image"

type Action interface {
	Do(img *image.NRGBA)
	Undo(img *image.NRGBA)
}