package main

import "strconv"

func formated_amount(amount float32) string {
	return strconv.FormatFloat(float64(amount), 'f', -1, 64)
}
