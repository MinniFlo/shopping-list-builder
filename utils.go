package main

import (
	"math"
	"strconv"
	"unicode/utf8"
)

func rightPadUnicodeConform(s string, pad_value int) string {
	pad_amt := pad_value - utf8.RuneCountInString(s)

	switch {
	case pad_amt == 0:
		return s
	case pad_amt > 0:
		runes := []rune(s)
		for range pad_amt {
			runes = append(runes, rune(' '))
		}
		return string(runes)
	case pad_amt < 0:
		runes := []rune(s)
		runes = append(runes[:pad_value-1], rune('…'))
		return string(runes)
	}

	return s
}

func RoundToThreeDigitsAfterPeriode(value float64) float64 {
	return math.Round((value)*1000.0) / 1000.0
}

func formatedAmount(amount float64) string {
	if amount == 0.0 {
		return ""
	}

	return strconv.FormatFloat(
		RoundToThreeDigitsAfterPeriode(amount),
		'f',
		-1,
		64,
	)
}

func idFromIndices(i, j int) int {
	return i*100 + j
}
