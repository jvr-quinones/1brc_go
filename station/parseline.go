package station

import (
	"errors"
	"strconv"
	"strings"
)

const sep = ";"

var ErrNoSep = errors.New("no separator found")

// Given string with a text an a number separated by a semicolon,
// the function will return two variables with the number parsed
// as a float and a string.
//
// If the number cannot be parsed, it will return an error.
func ParseLineFloat(str string) (string, float64, error) {
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

// Given string with a text an a number separated by a semicolon,
// the function will return two variables with the number parsed
// as an integer and a string.
//
// If the number cannot be parsed, it will return an error.
func ParseLineInt(str string) (string, int64, error) {
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
