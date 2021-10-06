package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strings"
)

const (
	BYTE = 1.0 << (10 * iota)
	KIBIBYTE
	MEBIBYTE
	GIBIBYTE
)

const ALPHA_NUM = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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

func UuidStrToBase64Str(uuidStr string) (string, error) {
	uuidObj, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", nil
	}
	uuidBytes, _ := uuidObj.MarshalBinary()
	return base64.RawURLEncoding.EncodeToString(uuidBytes), nil
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomAlphaNumStr(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = ALPHA_NUM[b%byte(len(ALPHA_NUM))]
	}
	return string(bytes)
}

func GetMyIpCountry() (string, string) {
	resp, err := http.Get("http://ip-api.com/json/?fields=query,country") // TODO: create own url for retrieving!
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", ""
	}
	return data["query"], data["country"]
}

func SplitArgsStr(argsStr string) []string {
	var args []string
	quoted := false
	startPos := 0
	n := len(argsStr)
	for i := 1; i < n; i++ {
		prev := argsStr[i-1]
		if prev == '"' {
			if quoted {
				args = append(args, argsStr[startPos:i-1])
				quoted = false
				startPos = i + 1
			} else {
				quoted = true
				startPos = i
			}
		} else if argsStr[i] == ' ' && !quoted {
			if prev != ' ' {
				args = append(args, argsStr[startPos:i])
				startPos = i
			}
			startPos++
		}
	}
	if n > 0 {
		if argsStr[n-1] == '"' {
			args = append(args, argsStr[startPos:n-1])
		} else {
			args = append(args, argsStr[startPos:n])
		}
	}
	return args
}
