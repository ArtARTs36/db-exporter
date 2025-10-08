package diagram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/imagedraw"
	"image/png"
)

type Creator struct {
	graphBuilder *GraphBuilder
	encoder      *png.Encoder
}

func NewCreator(
	graphBuilder *GraphBuilder,
) *Creator {
	return &Creator{
		graphBuilder: graphBuilder,
		encoder: &png.Encoder{
			CompressionLevel: png.NoCompression,
		},
	}
}

func (c *Creator) Create(
	ctx context.Context,
	tables *schema.TableMap,
	spec *config.DiagramExportSpec,
) ([]byte, error) {
	img, err := c.graphBuilder.Build(ctx, tables, spec)
	if err != nil {
		return nil, err
	}

	if spec.Style.Background.Grid != nil {
		img = imagedraw.AddBackground(
			img,
			imagedraw.GridFor(
				img,
				spec.Style.Background.Grid.CellSize,
				spec.Style.Background.Grid.LineColor.Color,
				spec.Style.Background.Color.Color,
			),
		)
	}

	var buf bytes.Buffer

	if err = c.encoder.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode to image: %w", err)
	}

	return buf.Bytes(), nil
}
