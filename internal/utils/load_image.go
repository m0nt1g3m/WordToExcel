package utils

import (
	"bytes"
	"fmt"
	"image/png"
	"wordtoexcel/internal/gui/assets"

	"gioui.org/op/paint"
)

func LoadPNG(name string) paint.ImageOp {
	fileName := fmt.Sprintf("public/%s", name)
	file, err := assets.PngImgs.ReadFile(fileName)
	if err != nil {
		return paint.ImageOp{}
	}
	image, err := png.Decode(bytes.NewReader(file))
	if err != nil {
		return paint.ImageOp{}
	}
	return paint.NewImageOp(image)
}
