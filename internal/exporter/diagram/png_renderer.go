package diagram

import (
	"context"
	"fmt"
	"github.com/flopp/go-findfont"
	"github.com/kanrichan/resvg-go"
	"log/slog"
	"os"
)

type PNGRenderer struct {
}

func NewPNGRenderer() *PNGRenderer {
	return &PNGRenderer{}
}

func (r *PNGRenderer) Render(
	ctx context.Context,
	svg []byte,
	fontName string,
) ([]byte, error) {
	worker, err := resvg.NewDefaultWorker(ctx)
	defer func() {
		if worker == nil {
			return
		}

		err = worker.Close()
		if err != nil {
			slog.Error("failed to close resvg worker", slog.Any("err", err))
		}
	}()
	if err != nil {
		return nil, err
	}

	fdb, err := worker.NewFontDBDefault()
	if err != nil {
		return nil, err
	}

	if err = r.loadFont(ctx, fdb, fontName); err != nil {
		return nil, fmt.Errorf("load font: %w", err)
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

func (r *PNGRenderer) loadFont(ctx context.Context, fdb *resvg.FontDB, font string) error {
	slog.DebugContext(ctx, "[png-renderer] loading font", slog.String("font.name", font))

	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", font))
	if err != nil {
		return fmt.Errorf("find font: %w", err)
	}

	slog.DebugContext(ctx, "[png-renderer] reading font", slog.String("font.path", fontPath))

	fnt, err := os.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("failed to open font: %w", err)
	}

	if err = fdb.LoadFontData(fnt); err != nil {
		return err
	}

	slog.DebugContext(ctx, "[png-renderer] font loaded to db", slog.String("font.path", fontPath))

	return nil
}
