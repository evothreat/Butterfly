package utils

import (
	"regexp"
)

var listStringRegex = regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)

func IsValidListString(str string) bool {
	return listStringRegex.MatchString(str)
}
