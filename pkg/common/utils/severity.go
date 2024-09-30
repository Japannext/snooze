package utils

import (
	"strings"
)

func GuessSeverityNumber(text string) int32 {
	switch strings.ToLower(text) {
		case "warning", "warn":
			return 13
		case "err", "error":
			return 17
		case "info", "ok":
			return 9
		case "debug":
			return 5
		case "trace":
			return 1
		case "fatal", "crit", "critical":
			return 21
	}
	return 16
}
