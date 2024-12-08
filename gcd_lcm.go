package utils

func GCD(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(integers ...int64) int64 {
	if len(integers) == 0 {
		return 0
	} else if len(integers) == 1 {
		return integers[0]
	}
	a := integers[0]
	b := integers[1]
	result := a * b / GCD(a, b)
	for i := 2; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}
