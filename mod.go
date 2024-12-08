package utils

func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}

func Mod64(a, b int64) int64 {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}
