package readline

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type ReadlineFloat struct {
	Name string
	Val  float64
}

type ReadlineInt struct {
	Name string
	Val  int64
}

const sep = ";"

var (
	ErrNoSep  = errors.New("no separator found")
	findRegex = regexp.MustCompile(`(;\d+)\.`)
)

func ReadAsFloat(str string) (*ReadlineFloat, error) {
	strName, strVal, hasSep := strings.Cut(str, sep)
	if !hasSep {
		return nil, ErrNoSep
	}

	val, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		return nil, err
	}

	return &ReadlineFloat{strName, val}, nil
}

func ReadAsInt(str string) (*ReadlineInt, error) {
	strName, strVal, hasSep := strings.Cut(
		findRegex.ReplaceAllString(str, `$1`),
		sep,
	)
	if !hasSep {
		return nil, ErrNoSep
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return nil, err
	}

	return &ReadlineInt{strName, int64(val)}, nil
}
