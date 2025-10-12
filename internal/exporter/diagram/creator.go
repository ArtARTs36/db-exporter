package diagram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/templates"
	"github.com/kanrichan/resvg-go"
	"io/ioutil"
	"log/slog"
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
	spec *config.DiagramExportSpec,
) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := c.graphBuilder.Build(ctx, tables, spec, buf)
	if err != nil {
		return nil, err
	}

	if spec.Style.Background.Grid != nil {
		buf = bytes.NewBuffer(c.injectGrid(buf.Bytes(), c.buildGridString(spec)))
	}

	bb, err := c.renderer.Render(ctx, buf.Bytes(), spec.Style.Font.Family)
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func (c *Creator) render(svg []byte) ([]byte, error) {
	worker, err := resvg.NewDefaultWorker(context.Background())
	defer func() {
		if worker == nil {
			return
		}

		err = worker.Close()
		slog.Error("failed to close resvg worker", slog.Any("err", err))
	}()
	if err != nil {
		return nil, err
	}

	fdb, err := worker.NewFontDBDefault()
	if err != nil {
		return nil, err
	}

	fnt, err := templates.FS.Open("diagram/arialmt.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to open font: %w", err)
	}

	fntData, err := ioutil.ReadAll(fnt)
	if err != nil {
		return nil, fmt.Errorf("failed to read font: %w", err)
	}

	if err = fdb.LoadFontData(fntData); err != nil {
		return nil, err
	}

	tree, err := worker.NewTreeFromData(svg, &resvg.Options{
		Dpi: 96.0,

		// Set rendering modes
		ShapeRenderingMode: resvg.ShapeRenderingModeGeometricPrecision,
		TextRenderingMode:  resvg.TextRenderingModeOptimizeLegibility,
		ImageRenderingMode: resvg.ImageRenderingModeOptimizeQuality,
	})
	defer func() {
		if tree == nil {
			return
		}

		err = tree.Close()
		if err != nil {
			slog.Error("failed to close tree", slog.Any("err", err))
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("create tree: %w", err)
	}

	width, height, err := tree.GetSize()
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	pixmap, err := worker.NewPixmap(uint32(width), uint32(height))
	defer func() {
		if pixmap == nil {
			return
		}

		err = pixmap.Close()
		if err != nil {
			slog.Error("failed to close pixmap", slog.Any("err", err))
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("create pixmap: %w", err)
	}

	if err = tree.ConvertText(fdb); err != nil {
		return nil, fmt.Errorf("convert text: %w", err)
	}

	if err = tree.Render(resvg.TransformIdentity(), pixmap); err != nil {
		return nil, fmt.Errorf("render tree: %w", err)
	}

	return pixmap.EncodePNG()
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
