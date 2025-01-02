package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type InputFileReader interface {
	io.ReadCloser
	ReadAll() ([]byte, error)
	ReadLine() (line []byte, isPrefix bool, err error)
	ReadLines(callback func(string) error) error
}

type inputFileReader struct {
	file      *os.File
	bufreader *bufio.Reader
}

func OpenInputFile(name string) (InputFileReader, error) {
	filename := fmt.Sprintf("./inputs/%s", name)
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	ifr := inputFileReader{
		file:      file,
		bufreader: nil,
	}
	return &ifr, nil
}

func OpenAndReadAll(name string) ([]byte, error) {
	ifr, err := OpenInputFile(name)
	if err != nil {
		return nil, err
	}
	data, err := ifr.ReadAll()
	if err != nil {
		return nil, err
	}
	err = ifr.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func OpenAndReadLines(name string, callback func(string) error) error {
	ifr, err := OpenInputFile(name)
	if err != nil {
		return err
	}
	err = ifr.ReadLines(callback)
	if err != nil {
		return err
	}
	err = ifr.Close()
	if err != nil {
		return err
	}
	return nil
}

func OpenAndReadRegex(name string, regex string, allMustMatch bool) ([][]string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	var results [][]string
	err = OpenAndReadLines(name, func(s string) error {
		m := re.FindStringSubmatch(s)
		if m == nil {
			if allMustMatch {
				return fmt.Errorf("non-matching line: %s", s)
			}
		} else {
			results = append(results, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type MultiRegex struct {
	Regex     string
	MatchFunc func([]string) error
	re        *regexp.Regexp
}

func OpenAndReadMultipleRegex(name string, regexes []MultiRegex, allMustMatch bool) error {
	for i := range regexes {
		var err error
		regexes[i].re, err = regexp.Compile(regexes[i].Regex)
		if err != nil {
			return err
		}
	}
	err := OpenAndReadLines(name, func(s string) error {
		for _, rx := range regexes {
			m := rx.re.FindStringSubmatch(s)
			if m != nil {
				err := rx.MatchFunc(m)
				if err != nil {
					return err
				}
				return nil
			}
		}
		if allMustMatch {
			return fmt.Errorf("line failed to match any regexes")
		}
		return nil
	})
	return err
}

func StringsToInts(s []string, positions ...int) ([]int, error) {
	var results []int
	for _, p := range positions {
		v, err := strconv.Atoi(s[p])
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

func (ifr *inputFileReader) Read(p []byte) (n int, err error) {
	return ifr.bufreader.Read(p)
}

func (ifr *inputFileReader) Close() error {
	return ifr.file.Close()
}

func (ifr *inputFileReader) ReadLine() (line []byte, isPrefix bool, err error) {
	if ifr.bufreader == nil {
		ifr.bufreader = bufio.NewReader(ifr.file)
	}
	return ifr.bufreader.ReadLine()
}

func (ifr *inputFileReader) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(ifr.file)
}

func (ifr *inputFileReader) ReadLines(callback func(string) error) error {
	scanner := bufio.NewScanner(ifr.file)
	var err error
	for scanner.Scan() {
		err = callback(scanner.Text())
		if err != nil {
			break
		}
	}
	if err == nil {
		err = scanner.Err()
	}
	return err
}
