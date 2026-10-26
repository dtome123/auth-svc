package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRegex = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)(d|h|ms|s|m|us|µs|ns)`)

func ParseFlexibleDuration(input string) (time.Duration, error) {
	matches := durationRegex.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	var total time.Duration
	for _, match := range matches {
		valueStr, unit := match[1], strings.ToLower(match[2])
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %s: %v", valueStr, err)
		}

		var dur time.Duration
		switch unit {
		case "d":
			dur = time.Duration(value * float64(24*time.Hour))
		case "h":
			dur = time.Duration(value * float64(time.Hour))
		case "m":
			dur = time.Duration(value * float64(time.Minute))
		case "s":
			dur = time.Duration(value * float64(time.Second))
		case "ms":
			dur = time.Duration(value * float64(time.Millisecond))
		case "us", "µs":
			dur = time.Duration(value * float64(time.Microsecond))
		case "ns":
			dur = time.Duration(value * float64(time.Nanosecond))
		default:
			return 0, fmt.Errorf("unsupported unit: %s", unit)
		}

		total += dur
	}

	return total, nil
}
