package controllers

import (
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

var listStringRegex = regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)

func ValuesDictToWhere(dict echo.Map) (string, []interface{}) {
	cols := ""
	vals := make([]interface{}, 0, 15)
	for k, v := range dict {
		cols += k + "=?,"
		vals = append(vals, v)
	}
	cols = strings.TrimSuffix(cols, ",")
	return cols, vals
}

func IsValidListString(str string) bool {
	return listStringRegex.MatchString(str)
}
