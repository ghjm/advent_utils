package utils

import "os"

func MustWriteString(f *os.File, s string) {
	_, err := f.WriteString(s)
	if err != nil {
		panic(err)
	}
}
