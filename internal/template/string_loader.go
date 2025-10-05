package template

import (
	"fmt"
	"strings"
	"time"

	"github.com/tyler-sommer/stick"
)

type StringLoader struct{}

func NewStringLoader() *StringLoader {
	return &StringLoader{}
}

func (l *StringLoader) Load(value string) (stick.Template, error) {
	return &fileTemplate{
		name:   fmt.Sprintf("string-%d", time.Now().UnixMilli()),
		reader: strings.NewReader(value),
	}, nil
}
