package diagram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/imageencoder"
)

type Creator struct {
	graphBuilder   *GraphBuilder
	encoderManager *imageencoder.Manager
}

func NewCreator(
	graphBuilder *GraphBuilder,
	encoderManager *imageencoder.Manager,
) *Creator {
	return &Creator{
		graphBuilder:   graphBuilder,
		encoderManager: encoderManager,
	}
}

func (c *Creator) Create(
	ctx context.Context,
	tables *schema.TableMap,
	spec *config.DiagramExportSpec,
) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := c.graphBuilder.Build(ctx, tables, spec, buf)
	if err != nil {
		return nil, err
	}

	if spec.Style.Background.Grid != nil {
		return c.injectGrid(buf.Bytes(), c.buildGridString(spec)), nil
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

			injector = func(i int, char byte) {
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

func (c *Creator) buildGridString(spec *config.DiagramExportSpec) string {
	return fmt.Sprintf(`<defs>
    <pattern id="gridPattern" width="20" height="20" patternUnits="userSpaceOnUse">
      <!-- Horizontal lines -->
      <path d="M 0 0 L 20 0 M 0 20 L 20 20" stroke="%s" stroke-width="0.5"/>
      <!-- Vertical lines -->
      <path d="M 0 0 L 0 20 M 20 0 L 20 20" stroke="%s" stroke-width="0.5"/>
    </pattern>
  </defs>
<rect x="0" y="0" width="100%%" height="100%%" fill="%s"/>
<rect x="0" y="0" width="100%%" height="100%%" fill="url(#gridPattern)"/>
`,
		spec.Style.Background.Grid.LineColor.Hex(),
		spec.Style.Background.Grid.LineColor.Hex(),
		spec.Style.Background.Color.Hex(),
	)
}
