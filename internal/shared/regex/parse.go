package regex

import (
	"regexp"
)

func ParseSingleValue(exp *regexp.Regexp, str string) string {
	const needMatchesCount = 2

	matches := exp.FindStringSubmatch(str)
	if len(matches) < needMatchesCount {
		return ""
	}

	return matches[1]
}
