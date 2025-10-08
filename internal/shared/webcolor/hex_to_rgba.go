package webcolor

import (
	"image/color"
	"strconv"
)

func Hex(hex string) color.RGBA {
	values, _ := strconv.ParseUint(hex[1:], 16, 32)

	return color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
		A: 255,
	}
}
