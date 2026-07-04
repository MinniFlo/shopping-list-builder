package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type list_entry struct {
	id       int
	name     string
	amount   float64
	unit     unit
	category category
	state    stageing_state
}

func createListEntryFromString(s string, mapping map[string]int, not_staged bool) (list_entry, error) {
	re := regexp.MustCompile(`- \[(.)\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
	entry_match := re.FindStringSubmatch(s)

	if entry_match != nil {
		name := "ENTRY_MISSING"
		amount := 0.0
		unit := None
		category := UNDEFINED
		state := STAGED

		checked := entry_match[1] != " "
		amount_match := strings.Replace(entry_match[2], ",", ".", 1)
		unit_match := entry_match[3]
		name_match := entry_match[4]

		if value, err := strconv.ParseFloat(amount_match, 32); err == nil {
			amount = RoundToThreeDigitsAfterPeriode(value)
		}

		if len(name_match) > 0 {
			name = strings.TrimSpace(name_match)
		}

		if len(unit_match) > 0 {
			unit_string := strings.TrimSpace(strings.ToLower(unit_match))
			unit = UnitFromString(unit_string)
		}

		if category_int, ok := mapping[name]; ok {
			category = CategoryFromInt(category_int)
		}

		if checked || not_staged {
			state = NOT_STAGED
		}

		return list_entry{name: name, amount: amount, unit: unit, category: category, state: state}, nil
	}

	return list_entry{}, errors.New("Invalid list entry string!")
}

func (m *list_model) GetCurrentPositionID() int {
	return idFromIndices(m.collection_index, m.entry_index)
}
