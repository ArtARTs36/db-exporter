package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Parser func(content []byte) (*Config, error)

func YAMLParser() Parser {
	return unmarshalParser(yaml.Unmarshal)
}

func JSONParser() Parser {
	return unmarshalParser(json.Unmarshal)
}

func unmarshalParser(unmarshal func(data []byte, v any) error) Parser {
	return func(content []byte) (*Config, error) {
		var cfg Config
		err := unmarshal(content, &cfg)
		return &cfg, err
	}
}
