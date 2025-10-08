package webcolor

import (
	"fmt"
	"strings"
)

const shortHexLen = 4

func Fix(color string) string {
	if strings.HasPrefix(color, "#") && len(color) == shortHexLen {
		return fmt.Sprintf("%s%s", color, color[1:shortHexLen])
	}

	return color
}
