package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type imageRequest struct {
	RequestPath    string
	SourceFilePath string
	CachedFilePath string
	Size           uint
}

func (r *imageRequest) resize(watermark image.Image) error {
	var scaledImage image.Image
	file, err := os.Open(r.SourceFilePath)

	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)

	if err != nil {
		return err
	}
	defer file.Close()

	if watermark != nil {
		imgBounds := img.Bounds()
		wtmBounds := watermark.Bounds()

		newImage := image.NewRGBA(imgBounds)

		// Watermark: http://stackoverflow.com/questions/16100023/manipulating-watermark-images-with-go
		offset := image.Pt(15, imgBounds.Dy()-wtmBounds.Dy()-15)

		draw.Draw(newImage, imgBounds, img, image.ZP, draw.Src)
		draw.Draw(newImage, wtmBounds.Add(offset), watermark, image.ZP, draw.Over)

		// m := resize.Resize(256, 0, img, resize.Lanczos3)
		scaledImage = resize.Thumbnail(r.Size, r.Size, newImage, resize.Lanczos3)
	} else {
		scaledImage = resize.Thumbnail(r.Size, r.Size, img, resize.Lanczos3)
	}

	dstDir := filepath.Dir(r.CachedFilePath)

	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, os.ModeDir|os.ModePerm)
	}

	// possibly race condition with another routine
	if _, err := os.Stat(r.CachedFilePath); os.IsExist(err) {
		return nil
	}

	out, err := os.Create(r.CachedFilePath)

	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, scaledImage, nil)
}

func getWatermark(path string) (image.Image, error) {
	var img image.Image
	var err error

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err = png.Decode(file)

	if err != nil {
		return nil, err
	}

	return img, nil
}
