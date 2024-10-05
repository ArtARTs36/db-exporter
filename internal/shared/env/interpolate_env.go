package env

import "os"

type interpolateEnv struct {
}

func (*interpolateEnv) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}
