package common

import (
	"strconv"
	"strings"
)

func HexToDecimal(hex string) (int64, error) {
	hex = strings.TrimPrefix(strings.TrimPrefix(hex, "0x"), "0X")
	num, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}
