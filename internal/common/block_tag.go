package common

func IsBlockTag(tag string) bool {
	tags := map[string]struct{}{
		"earliest":  {},
		"latest":    {},
		"pending":   {},
		"safe":      {},
		"finalized": {},
	}

	_, isTag := tags[tag]

	if isTag {
		return true
	}

	_, err := HexToDecimal(tag)

	if err != nil {
		return false
	}

	return true
}
