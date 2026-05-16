package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type list_entry struct {
	name     string
	amount   float32
	unit     string
	category category
	staged   stageing_state
}

func (e *list_entry) formated_amount(amount float32) string {
	if (amount == 1.0 && e.unit == "") {
		return ""
	}
	return strconv.FormatFloat(float64(amount), 'f', -1, 64)
}

func createListEntryFromString(s string, mapping map[string]int) (list_entry, error) {
	re := regexp.MustCompile(`- \[.\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
	incredient_match := re.FindStringSubmatch(s)

	if incredient_match != nil {
		name := "INCREDIENT_MISSING"
		amount := 1.0
		unit := ""
		category := UNDEFINED

		if value, err := strconv.ParseFloat(strings.Replace(incredient_match[1], ",", ".", 1), 32); err == nil {
			amount = value
		}

		if len(incredient_match[3]) > 0 {
			name = strings.TrimSpace(incredient_match[3])
		}

		if len(incredient_match[2]) > 0 {
			unit = strings.TrimSpace(incredient_match[2])
		}

		if category_int, ok := mapping[name]; ok {
			category = CategoryFromInt(category_int)
		}

		return list_entry{name: name, amount: float32(amount), unit: unit, category: category, staged: STAGED}, nil
	}

	return list_entry{}, errors.New("Invalid incredient string!")
}
