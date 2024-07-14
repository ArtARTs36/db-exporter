package fs

type Driver interface {
	Exists(path string) bool
	ReadFile(path string) ([]byte, error)
	Mkdir(path string) error
	Write(path string, content []byte) (FileInfo, error)
}

type FileInfo struct {
	Path string
	Size int64
}
