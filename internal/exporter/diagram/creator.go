package diagram

import (
	"bytes"
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Creator struct {
	graphBuilder *GraphBuilder
	renderer     *PNGRenderer
}

func NewCreator(
	graphBuilder *GraphBuilder,
) *Creator {
	return &Creator{
		graphBuilder: graphBuilder,
		renderer:     NewPNGRenderer(),
	}
}

func (c *Creator) Create(
	ctx context.Context,
	tables *schema.TableMap,
	spec *Specification,
) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := c.graphBuilder.Build(ctx, tables, spec, buf)
	if err != nil {
		return nil, err
	}

	if spec.Style.Background.Grid != nil {
		buf = bytes.NewBuffer(c.injectGrid(buf.Bytes(), c.buildGridString(spec)))
	}

	if spec.Image.Format == ImageFormatPNG {
		var rerr error

		bb, rerr := c.renderer.Render(ctx, buf.Bytes(), spec.Style.Font.Family)
		if rerr != nil {
			return nil, rerr
		}

		buf = bytes.NewBuffer(bb)
	}

	return buf.Bytes(), nil
}

func (c *Creator) injectGrid(content []byte, grid string) []byte {
	buf := bytes.NewBuffer([]byte{})

	inRoot := false

	var injector func(i int, char byte)

	injector = func(i int, char byte) {
		if char == '<' && i < len(content)-4 && content[i+1] == 's' && content[i+2] == 'v' && content[i+3] == 'g' {
			inRoot = true
			buf.WriteByte(char)
			return
		} else if inRoot && char == '>' {
			buf.WriteByte(char)
			buf.WriteRune('\n')
			buf.WriteString(grid)

			injector = func(_ int, char byte) {
				buf.WriteByte(char)
			}

			inRoot = false
			return
		}

		buf.WriteByte(char)
	}

	for i, char := range content {
		injector(i, char)
	}

	return buf.Bytes()
}

func (c *Creator) buildGridString(spec *Specification) string {
	return fmt.Sprintf(`<defs>
    <pattern id="gridPattern" width="%d" height="%d" patternUnits="userSpaceOnUse">
      <!-- Horizontal lines -->
      <path d="M 0 0 L %d 0 M 0 %d L %d %d" stroke="%s" stroke-width="0.5"/>
      <!-- Vertical lines -->
      <path d="M 0 0 L 0 %d M %d 0 L %d %d" stroke="%s" stroke-width="0.5"/>
    </pattern>
  </defs>
<rect x="0" y="0" width="100%%" height="100%%" fill="%s"/>
<rect x="0" y="0" width="100%%" height="100%%" fill="url(#gridPattern)"/>
`,
		// pattern width and height
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,

		// horizontal path
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,

		spec.Style.Background.Grid.LineColor.Hex(),

		// vertical path
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.CellSize,
		spec.Style.Background.Grid.LineColor.Hex(),

		// background color
		spec.Style.Background.Color.Hex(),
	)
}
