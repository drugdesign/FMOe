package cpf

import (
	"fmt"
	"strconv"
	"strings"
)

// StringOutOfRange error
type StringOutOfRange struct {
	String string
	Index  int
}

func (err *StringOutOfRange) Error() string {
	return fmt.Sprintf("string out of range: %v of %v", err.Index, err.String)
}

// IsStringOutOfRange checks error isb StringOutOfRange or not
func IsStringOutOfRange(e error) bool {
	_, ok := e.(*StringOutOfRange)
	return ok
}

func slice(s string, start, end int) (string, error) {
	if start < 0 {
		return "", &StringOutOfRange{String: s, Index: start}
	}
	if len(s) < end {
		return "", &StringOutOfRange{String: s, Index: end}
	}
	return s[start:end], nil
}

func sliceTrim(s string, start, end int) (string, error) {
	v, err := slice(s, start, end)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(v), nil
}

func intField(s string, start, end int) (int, error) {
	s, err := sliceTrim(s, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func floatField(s string, start, end int) (float64, error) {
	s, err := sliceTrim(s, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}
