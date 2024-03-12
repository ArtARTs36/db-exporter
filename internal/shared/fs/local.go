package fs

import (
	"os"
	"syscall"
	"time"
)

type Local struct {
}

func NewLocal() *Local {
	return &Local{}
}

func (*Local) Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func (*Local) Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

func (*Local) CreateFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0755)
}

func (*Local) Stat(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	sysStat := stat.Sys().(*syscall.Stat_t)
	ctime := time.Unix(sysStat.Birthtimespec.Sec, sysStat.Birthtimespec.Nsec)

	return &FileInfo{
		Path:      path,
		CreatedAt: ctime,
		UpdatedAt: stat.ModTime(),
	}, nil
}
