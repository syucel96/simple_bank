package util

import (
	"fmt"
	"math"
	"strconv"
)

func NewFloat(n int64, f float64) float64 {
	return math.Round((float64(n)+f)*100) / 100
}

func AddDecimals(d1, d2 string) (string, error) {
	f1, err := strconv.ParseFloat(d1, 64)
	if err != nil {
		return "First argument", err
	}
	f2, err := strconv.ParseFloat(d2, 64)
	if err != nil {
		return "Secound argument", err
	}
	return fmt.Sprintf("%.2f", f1+f2), nil
}

func SubtractDecimals(d1, d2 string) (string, error) {
	f1, err := strconv.ParseFloat(d1, 64)
	if err != nil {
		return "First argument", err
	}
	f2, err := strconv.ParseFloat(d2, 64)
	if err != nil {
		return "Secound argument", err
	}
	return fmt.Sprintf("%.2f", f1-f2), nil
}

func ChangeDecimalSign(d1 string) string {
	if d1[0] == '-' {
		return d1[1:]
	} else {
		return fmt.Sprintf("-%s", d1)
	}
}

func ParseFloat(d string) float64 {
	res, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0
	}
	return res
}
