package mysql

import "strings"

func ParseEnumType(definition string) []string {
	if definition == "" {
		return nil
	}

	values := []string{}

	for i := 0; i < len(definition); i++ {
		if definition[i] != '\'' {
			continue
		}

		value := strings.Builder{}

		for j := i + 1; j < len(definition) && definition[j] != '\''; j++ {
			_ = value.WriteByte(definition[j])
		}

		values = append(values, value.String())
		i += value.Len() + 1
	}

	return values
}
