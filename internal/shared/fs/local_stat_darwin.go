//go:build darwin

package fs

import (
	"os"
	"syscall"
	"time"
)

func (*Local) Stat(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	sysStat, err := stat.Sys().(*syscall.Stat_t)
	if err != nil {
		return nil, err
	}

	ctime := time.Unix(sysStat.Birthtimespec.Sec, sysStat.Birthtimespec.Nsec)

	return &FileInfo{
		Path:      path,
		CreatedAt: ctime,
		UpdatedAt: stat.ModTime(),
	}, nil
}
