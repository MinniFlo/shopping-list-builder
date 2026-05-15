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

func createListEntryFromString(s string) (list_entry, error) {
	re := regexp.MustCompile(`- \[.\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
	incredient_match := re.FindStringSubmatch(s)

	if incredient_match != nil {
		name := "INCREDIENT_MISSING"
		amount := 1.0
		unit := ""

		if value, err := strconv.ParseFloat(incredient_match[1], 32); err == nil {
			amount = value
		}

		if len(incredient_match[3]) > 0 {
			name = strings.TrimSpace(incredient_match[3])
		}

		if len(incredient_match[2]) > 0 {
			unit = strings.TrimSpace(incredient_match[2])
		}

		return list_entry{name: name, amount: float32(amount), unit: unit, category: UNDEFINED, staged: STAGED}, nil
	}

	return list_entry{}, errors.New("Invalid incredient string!")
}
