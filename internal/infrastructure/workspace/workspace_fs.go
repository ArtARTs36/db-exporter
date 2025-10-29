package workspace

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"log/slog"

	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type FSWorkspace struct {
	cfg Config

	fs fs.Driver

	writerFn func(ctx context.Context, path string, writer func(buffer iox.Writer) error) error

	store *store
}

type Config struct {
	Directory  string
	FilePrefix string
	SkipExists bool
}

func newFSWorkspace(
	cfg Config,
	fs fs.Driver,
	store *store,
) *FSWorkspace {
	ws := &FSWorkspace{
		cfg:   cfg,
		fs:    fs,
		store: store,
	}

	ws.setupWriter()

	return ws
}

func (w *FSWorkspace) Write(ctx context.Context, path string, writer func(buffer iox.Writer) error) error {
	return w.writerFn(ctx, path, writer)
}

func (w *FSWorkspace) setupWriter() {
	w.writerFn = w.write

	if w.cfg.SkipExists {
		w.writerFn = w.writeNewFile
	}
}

func (w *FSWorkspace) writeNewFile(ctx context.Context, path string, writer func(buffer iox.Writer) error) error {
	if w.fs.Exists(path) {
		return nil
	}

	return w.write(ctx, path, writer)
}

func (w *FSWorkspace) write(ctx context.Context, filename string, writer func(buffer iox.Writer) error) error {
	path := w.pathTo(filename)

	file, err := w.fs.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.WarnContext(ctx,
				"[workspace] failed to close file",
				slog.Any("file", file.Name()),
				slog.Any("err", cerr),
			)
		}
	}()

	buf := iox.NewWriter()

	err = writer(buf)
	if err != nil {
		return fmt.Errorf("write to buffer: %w", err)
	}

	size, err := file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	w.store.Add(&fs.FileInfo{
		Path: path,
		Size: int64(size),
	})

	return nil
}

func (w *FSWorkspace) pathTo(filename string) string {
	return fmt.Sprintf("%s/%s%s", w.cfg.Directory, w.cfg.FilePrefix, filename)
}
