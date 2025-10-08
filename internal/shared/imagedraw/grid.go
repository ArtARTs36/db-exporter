package imagedraw

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
)

func GridFor(img image.Image, spacing int, lineColor color.Color, bgColor color.Color) *image.RGBA {
	return Grid(
		img.Bounds().Dx(),
		img.Bounds().Dy(),
		spacing,
		lineColor,
		bgColor,
	)
}

func Grid(
	width, height, spacing int,
	lineColor color.Color,
	backgroundColor color.Color,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	for y := 0; y < height; y += spacing {
		for x := 0; x < width; x++ {
			img.Set(x, y, lineColor)
		}
	}

	for x := 0; x < width; x += spacing {
		for y := 0; y < height; y++ {
			img.Set(x, y, lineColor)
		}
	}

	return img
}

func AddBackground(object image.Image, background image.Image) (image.Image, error) {
	drawImage, ok := object.(draw.Image)
	if !ok {
		return object, errors.New("object does not implement draw.Image")
	}

	newImage := image.NewRGBA(image.Rect(0, 0, object.Bounds().Dx(), object.Bounds().Dy()))

	draw.Draw(newImage, object.Bounds(), background, image.Point{}, draw.Src)

	draw.Draw(newImage, object.Bounds(), drawImage, image.Point{}, draw.Over)

	return newImage, nil
}
