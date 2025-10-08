package imagedraw

import (
	"image"
	"image/color"
	"image/draw"
)

func GridFor(img image.Image, spacing int, lineColor color.RGBA) *image.RGBA {
	return Grid(
		img.Bounds().Dx(),
		img.Bounds().Dy(),
		spacing,
		lineColor,
	)
}

func Grid(
	width, height, spacing int,
	lineColor color.RGBA,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

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

func AddBackground(object image.Image, background image.Image) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, object.Bounds().Dx(), object.Bounds().Dy()))

	// Сначала размещаем изображение на задний план
	draw.Draw(newImage, object.Bounds(), background, image.Point{}, draw.Src)

	// Затем накладываем основное изображение поверх него
	draw.Draw(newImage, object.Bounds(), object.(draw.Image), image.Point{}, draw.Over)

	return newImage
}
