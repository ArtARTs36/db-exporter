package fs

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

func CreateFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0755)
}
