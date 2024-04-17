package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHexToDecimal(t *testing.T) {
	testCases := []struct {
		hex      string
		expected int64
		hasError bool
	}{
		{"1", 1, false},
		{"0x1", 1, false},
		{"0x242d504", 37934340, false},
		{"", 0, true},
		{"XXXX", 0, true},
	}

	for _, tc := range testCases {
		actual, err := HexToDecimal(tc.hex)

		if tc.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		}
	}
}
