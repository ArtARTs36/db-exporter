package imageencoder

import "fmt"

type Manager struct {
	encoders map[Format]Encoder
}

func NewManager() *Manager {
	return &Manager{
		encoders: map[Format]Encoder{ //nolint:exhaustive // not need
			FormatPNG:  &PNG{},
			FormatJPEG: &JPEG{},
		},
	}
}

func (m *Manager) For(format Format) (Encoder, error) {
	if enc, ok := m.encoders[format]; ok {
		return enc, nil
	}
	return nil, fmt.Errorf("unknown format: %s", format)
}
