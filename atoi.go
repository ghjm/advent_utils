package utils

import "strconv"

func MustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func MustAtoiHex(s string) int {
	v, err := strconv.ParseInt(s, 16, 0)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func MustAtoi64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func MustAtoiU64(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func MustAtoiHex64(s string) int64 {
	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func MustAtof(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}
