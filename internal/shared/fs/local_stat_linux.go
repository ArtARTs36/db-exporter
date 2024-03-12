//go:build linux

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

	ctime := time.Unix(int64(sysStat.Ctim.Sec), int64(sysStat.Ctim.Nsec))

	return &FileInfo{
		Path:      path,
		CreatedAt: ctime,
		UpdatedAt: stat.ModTime(),
	}, nil
}
