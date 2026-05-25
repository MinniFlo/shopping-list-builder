package main

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Style struct {
	black                             color.Color
	light_black                       color.Color
	white                             color.Color
	light_white                       color.Color
	primary                           color.Color
	highlight                         color.Color
	docStyle                          lipgloss.Style
	title                             lipgloss.Style
	table                             lipgloss.Style
	collection_header                 lipgloss.Style
	collection_header_selected        lipgloss.Style
	collection_header_amount          lipgloss.Style
	collection_header_amount_selected lipgloss.Style
	entry                             lipgloss.Style
	entry_odd                         lipgloss.Style
	entry_selected                    lipgloss.Style
	col                               lipgloss.Style
	name_col                          lipgloss.Style
	amount_col                        lipgloss.Style
	unit_col                          lipgloss.Style
}

func BuildStyle(bg_color, fg_color string) Style {

	black := lipgloss.Color(bg_color)
	light_black := lipgloss.Lighten(black, 0.05)
	white := lipgloss.Color(fg_color)
	light_white := lipgloss.Darken(white, 0.15)
	primary := lipgloss.Color("#ce82c1")
	highlight := lipgloss.Color("#8c8191")

	return Style{
		// Colors
		black:       black,
		light_black: light_black,

		white:       white,
		light_white: light_white,

		primary:   primary,
		highlight: highlight,

		docStyle: lipgloss.NewStyle().Padding(1, 2, 1, 2),

		title: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Italic(true).
			Foreground(black).
			Background(primary),

		// List/Table,
		table: lipgloss.NewStyle().MarginLeft(1),

		collection_header:                 lipgloss.NewStyle().Padding(0, 1, 0, 1).Bold(true),
		collection_header_selected:        lipgloss.NewStyle().Padding(0, 1, 0, 1).Bold(true).Foreground(black).Background(highlight),
		collection_header_amount:          lipgloss.NewStyle().Bold(true).Foreground(light_white).PaddingLeft(1),
		collection_header_amount_selected: lipgloss.NewStyle().Bold(true).Foreground(light_black).Background(highlight).PaddingLeft(1),

		entry:          lipgloss.NewStyle().Foreground(light_white),
		entry_odd:      lipgloss.NewStyle().Foreground(light_white).Background(light_black),
		entry_selected: lipgloss.NewStyle().Foreground(light_black).Background(highlight),

		col:        lipgloss.NewStyle().Padding(0, 1, 0, 1),
		name_col:   lipgloss.NewStyle().Width(35).Padding(0, 1, 0, 1),
		amount_col: lipgloss.NewStyle().Width(9).AlignHorizontal(lipgloss.Right).Padding(0, 1, 0, 1),
		unit_col:   lipgloss.NewStyle().Padding(0, 1, 0, 0),
	}
}
