package fs

type Driver interface {
	Exists(path string) bool
	Mkdir(path string) error
	CreateFile(path string, content []byte) error
}
