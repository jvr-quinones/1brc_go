package readline

import (
	"errors"
	"strconv"
	"strings"
)

const sep = ";"

var ErrNoSep = errors.New("no separator found")

func ReadAsFloat(str string) (string, float64, error) {
	strName, strVal, hasSep := strings.Cut(str, sep)
	if !hasSep {
		return "", 0.0, ErrNoSep
	}

	val, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return "", 0.0, err
	}

	return strName, val, nil
}

func ReadAsInt(str string) (string, int64, error) {
	strName, strVal, hasSep := strings.Cut(str, sep)
	if !hasSep {
		return "", 0, ErrNoSep
	}

	val, err := strconv.Atoi(strings.ReplaceAll(strVal, ".", ""))
	if err != nil {
		return "", 0, err
	}

	return strName, int64(val), nil
}
