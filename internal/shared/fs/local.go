package fs

import (
	"os"
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

func (*Local) Write(path string, content []byte) (FileInfo, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return FileInfo{}, err
	}

	size, err := f.Write(content)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}

	return FileInfo{Path: path, Size: int64(size)}, err
}
