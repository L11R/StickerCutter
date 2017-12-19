package main

import (
	"image"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"math"
)

func MakeSticker(img image.Image) image.Image {
	dc := gg.NewContext(512, 512)
	dc.DrawCircle(256, 256, 256)
	dc.Clip()

	var resizedImg image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		resizedImg = resize.Resize(0, 512, img, resize.Lanczos3)
		dc.DrawImage(resizedImg, -int(math.Floor(float64(resizedImg.Bounds().Dx() - resizedImg.Bounds().Dy()) / 2)), 0)
	} else if img.Bounds().Dy() > img.Bounds().Dx() {
		resizedImg = resize.Resize(512, 0, img, resize.Lanczos3)
		dc.DrawImage(resizedImg, 0, -int(math.Floor(float64(resizedImg.Bounds().Dy() - resizedImg.Bounds().Dx()) / 2)))
	} else {
		resizedImg = resize.Resize(512, 0, img, resize.Lanczos3)
		dc.DrawImage(resizedImg, 0, 0)
	}

	return dc.Image()
}