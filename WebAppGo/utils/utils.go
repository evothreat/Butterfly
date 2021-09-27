package utils

import (
	"crypto/rand"
	"encoding/hex"
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

func RandomBytes(n int) []byte {
	bs := make([]byte, n)
	rand.Read(bs)
	return bs
}

func RandomHexString(n int) string {
	return hex.EncodeToString(RandomBytes(n))
}
