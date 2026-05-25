package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type list_entry struct {
	name     string
	amount   float64
	unit     unit
	category category
	state    stageing_state
}

func (e *list_entry) formated_amount(amount float64) string {
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

func createListEntryFromString(s string, mapping map[string]int) (list_entry, error) {
	re := regexp.MustCompile(`- \[.\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
	incredient_match := re.FindStringSubmatch(s)

	if incredient_match != nil {
		name := "INCREDIENT_MISSING"
		amount := 0.0
		unit := None
		category := UNDEFINED

		amount_match := strings.Replace(incredient_match[1], ",", ".", 1)
		unit_match := incredient_match[2]
		name_match := incredient_match[3]

		if value, err := strconv.ParseFloat(amount_match, 32); err == nil {
			amount = RoundToThreeDigitsAfterPeriode(value)
		}

		if len(name_match) > 0 {
			name = strings.TrimSpace(incredient_match[3])
		}

		if len(unit_match) > 0 {
			unit_string := strings.TrimSpace(strings.ToLower(incredient_match[2]))
			unit = UnitFromString(unit_string)
		}

		if category_int, ok := mapping[name]; ok {
			category = CategoryFromInt(category_int)
		}

		return list_entry{name: name, amount: amount, unit: unit, category: category, state: STAGED}, nil
	}

	return list_entry{}, errors.New("Invalid incredient string!")
}
