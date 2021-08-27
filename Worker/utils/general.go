package utils

import (
	"fmt"
	"strings"
)

const (
	BYTE = 1.0 << (10 * iota)
	KIBIBYTE
	MEBIBYTE
	GIBIBYTE
)

func ToReadableSize(bytes uint64) string {
	unit := ""
	val := float32(bytes)
	switch {
	case bytes >= GIBIBYTE:
		unit = "GiB"
		val = val / GIBIBYTE
	case bytes >= MEBIBYTE:
		unit = "MiB"
		val = val / MEBIBYTE
	case bytes >= KIBIBYTE:
		unit = "KiB"
		val = val / KIBIBYTE
	case bytes >= BYTE:
		unit = "bytes"
	case bytes == 0:
		return "0"
	}
	strVal := strings.TrimSuffix(
		fmt.Sprintf("%.2f", val), ".00",
	)
	return fmt.Sprintf("%s %s", strVal, unit)
}
