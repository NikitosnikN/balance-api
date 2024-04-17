package common

func IsEthAddress(address string) bool {
	return len(address) == 42 && (address[:2] == "0x" || address[:2] == "0X")
}
