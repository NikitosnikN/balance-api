package nodepool

import "errors"

func determineTrueBalance(records map[string]int) (string, error) {
	if len(records) == 0 {
		return "", errors.New("empty records")
	}

	maxValue := 0
	maxKey := ""
	isDuplicate := false

	for key, value := range records {
		if value > maxValue {
			maxValue = value
			maxKey = key
			isDuplicate = false
		} else if value == maxValue {
			isDuplicate = true
		}
	}

	if isDuplicate {
		return "", errors.New("cannot determine true balance because of multiple highest values")
	}

	return maxKey, nil
}
