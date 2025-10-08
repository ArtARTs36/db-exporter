package diagram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/imagedraw"
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
	img, err := c.graphBuilder.Build(ctx, tables, spec)
	if err != nil {
		return nil, err
	}

	if spec.Style.Background.Grid != nil {
		img, err = imagedraw.AddBackground(
			img,
			imagedraw.GridFor(
				img,
				spec.Style.Background.Grid.CellSize,
				spec.Style.Background.Grid.LineColor.Color,
				spec.Style.Background.Color.Color,
			),
		)
		if err != nil {
			return nil, fmt.Errorf("add background: %w", err)
		}
	}

	var buf bytes.Buffer

	encoder, err := c.encoderManager.For(spec.Image.Format)
	if err != nil {
		return nil, err
	}

	if err = encoder.Encode(&buf, img, imageencoder.Options{
		CompressionLevel: spec.Image.Compression,
	}); err != nil {
		return nil, fmt.Errorf("encode to image: %w", err)
	}

	return buf.Bytes(), nil
}
