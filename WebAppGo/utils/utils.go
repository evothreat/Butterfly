package utils

import (
	"crypto/rand"
	"regexp"
)

var listStringRegex = regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)
var fileNameRegex = regexp.MustCompile(`^[\w\-. ]+$`)

func IsValidListString(str string) bool {
	return listStringRegex.MatchString(str)
}

func IsValidFilename(str string) bool {
	return fileNameRegex.MatchString(str)
}

func GetRandomBytes(n int) []byte {
	bs := make([]byte, n)
	rand.Read(bs)
	return bs
}
