package util

import (
	"fmt"
	"testing"
)

func TestParseXYResString(t *testing.T) {
	for _, tt := range []struct{
		res string
		expectX int
		expectY int
		expectErr error
	}{
		{"1024x768", 1024, 768, nil},
		{"1x1", 1, 1, nil},
		{"AxB", 0, 0, fmt.Errorf("failed to parse x value of resolution 'axb': strconv.Atoi: parsing \"a\": invalid syntax")},
		{"1xB", 0, 0, fmt.Errorf("failed to parse y value of resolution '1xb': strconv.Atoi: parsing \"b\": invalid syntax")},
		{"2x-1", 0, 0, fmt.Errorf("y value of resolution '2x-1' should be > 0")},
		{"-10x1", 0, 0, fmt.Errorf("x value of resolution '-10x1' should be > 0")},
	} {
		t.Run(tt.res, func(t *testing.T) {
			x, y, err := ParseXYResString(tt.res)
			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("error = nil, wanted %#v", tt.expectErr)
				} else if err.Error() != tt.expectErr.Error() {
					t.Errorf("error = %#v, wanted %#v", err, tt.expectErr)
				}
			} else if err != nil {
				t.Errorf("error = %#v, wanted nil", err)
			}
			if x != tt.expectX {
				t.Errorf("x = %#v, wanted %#v", x, tt.expectX)
			}
			if y != tt.expectY {
				t.Errorf("y = %#v, wanted %#v", y, tt.expectY)
			}
		})
	}
}
