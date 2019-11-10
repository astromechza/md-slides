package util

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseXYResString parses a given string as an AxB resolution.
// For example:
//     ParseXYRes("1024x768") -> 1024, 768, nil
func ParseXYResString(i string) (int, int, error) {
	i = strings.TrimSpace(strings.ToLower(i))
	parts := strings.Split(i, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("resolution '%s' did not contain one 'x'", i)
	}
	xRes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse x value of resolution '%s': %s", i, err)
	}
	yRes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse y value of resolution '%s': %s", i, err)
	}
	if xRes <= 0 {
		return 0, 0, fmt.Errorf("x value of resolution '%s' should be > 0", i)
	}
	if yRes <= 0 {
		return 0, 0, fmt.Errorf("y value of resolution '%s' should be > 0", i)
	}
	return int(xRes), int(yRes), nil
}
