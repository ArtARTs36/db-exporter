package fs

import "os"

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
