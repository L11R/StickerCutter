package main

import (
	"image"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"math"
	"io"
	"os/exec"
	"io/ioutil"
	"fmt"
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

func RemoveAudio(video io.Reader) error {
	// "-movflags", "frag_keyframe+empty_moov", "-f", "mp4",
	cmd := exec.Command("ffmpeg", "-y", "-i", "temp1.mp4", "-vcodec", "copy", "-an", "temp2.mp4")
	cmd.Stdin = video

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	outerr, _ := ioutil.ReadAll(stderr)
	fmt.Println(string(outerr))

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

