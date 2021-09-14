package utils

import (
	"regexp"
	"strings"
)

var listStringRegex = regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)

func ValuesMapToWhere(valuesMap map[string]interface{}) (string, []interface{}) {
	cols := ""
	vals := make([]interface{}, 0, 15)
	for k, v := range valuesMap {
		cols += k + "=?,"
		vals = append(vals, v)
	}
	cols = strings.TrimSuffix(cols, ",")
	return cols, vals
}

func IsValidListString(str string) bool {
	return listStringRegex.MatchString(str)
}
