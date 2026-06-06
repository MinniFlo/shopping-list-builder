package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

type list_entry_collection struct {
	name    string
	entries []list_entry
	amount  int
}

func (collection *list_entry_collection) buildHeader(selected bool, style ListStyle) string {
	var name_style lipgloss.Style
	var amount_style lipgloss.Style

	if selected {
		name_style = style.collection_header_selected
		amount_style = style.collection_header_amount_selected
	} else {
		name_style = style.collection_header
		amount_style = style.collection_header_amount
	}

	amount := lipgloss.Sprintf("x%d", collection.amount)
	return lipgloss.Sprintf("%s%s", name_style.Render(collection.name), amount_style.Render(amount))
}

func (collection *list_entry_collection) asTableRows(style ListStyle) [][]string {
	var tableData [][]string

	for _, entry := range collection.entries {
		name := style.name_col.Render(entry.name)
		amount := style.amount_col.Render(formatedAmount(entry.amount * float64(collection.amount)))
		unit := style.unit_col.Render(entry.unit.String())
		category := style.col.Render(entry.category.Symbol())
		state := style.col.Render(entry.state.String())

		tableData = append(
			tableData,
			[]string{name, amount, unit, category, state},
		)
	}

	return tableData
}

func (collection *list_entry_collection) buildTable(selected_row_index int, style ListStyle) string {
	rows := collection.asTableRows(style)

	t := table.New().
		BaseStyle(lipgloss.NewStyle().Foreground(style.colors.light_white).Background(style.colors.black)).
		BorderTop(false).
		BorderBottom(false).
		BorderColumn(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle()
			if collection.entries[row].state == NOT_STAGED {
				if col == 0 {
					s = s.Strikethrough(true).StrikethroughSpaces(false)
				}
				s = s.Faint(true)
			}
			switch {
			case row == selected_row_index:
				return s.Inherit(style.entry_selected)
			case row%2 == 0:
				return s.Inherit(style.entry)
			default:
				return s.Inherit(style.entry_odd)
			}
		}).
		Rows(rows...)

	return t.Render()
}
