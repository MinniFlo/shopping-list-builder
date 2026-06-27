package main

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Colors struct {
	black       color.Color
	light_black color.Color
	white       color.Color
	light_white color.Color
	primary     color.Color
	highlight   color.Color
}

type ListStyle struct {
	colors                            Colors
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

type ExportStyle struct {
	colors         Colors
	list           lipgloss.Style
	heading        lipgloss.Style
	entry          lipgloss.Style
	selected_entry lipgloss.Style
}

type ConfirmMenuStyle struct {
	colors        Colors
	minWidth      int
	minHeight     int
	menu          lipgloss.Style
	description   lipgloss.Style
	confirmButton lipgloss.Style
	cancelButton  lipgloss.Style
}

type Style struct {
	colors           Colors
	docStyle         lipgloss.Style
	title            lipgloss.Style
	listStyle        ListStyle
	exportStyle      ExportStyle
	confirmMenuStyle ConfirmMenuStyle
}

func BuildStyle(bg_color, fg_color string) Style {

	black := lipgloss.Color(bg_color)
	light_black := lipgloss.Lighten(black, 0.05)
	white := lipgloss.Color(fg_color)
	light_white := lipgloss.Darken(white, 0.15)
	primary := lipgloss.Color("#ce82c1")
	highlight := lipgloss.Color("#b2a7b7")

	colors := Colors{
		black:       black,
		light_black: light_black,
		white:       white,
		light_white: light_white,
		primary:     primary,
		highlight:   highlight,
	}

	listStyle := ListStyle{
		colors: colors,

		table: lipgloss.NewStyle().Margin(0, 1, 0, 1),

		collection_header:                 lipgloss.NewStyle().Padding(0, 1).Bold(true),
		collection_header_selected:        lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(black).Background(highlight),
		collection_header_amount:          lipgloss.NewStyle().Bold(true).Foreground(light_white).PaddingRight(1),
		collection_header_amount_selected: lipgloss.NewStyle().Bold(true).Foreground(light_black).Background(highlight).PaddingRight(1),

		entry:          lipgloss.NewStyle().Foreground(light_white),
		entry_odd:      lipgloss.NewStyle().Foreground(light_white).Background(light_black),
		entry_selected: lipgloss.NewStyle().Foreground(light_black).Background(highlight),

		col:        lipgloss.NewStyle().Padding(0, 1),
		name_col:   lipgloss.NewStyle().Width(35).Padding(0, 1),
		amount_col: lipgloss.NewStyle().Width(9).AlignHorizontal(lipgloss.Right).Padding(0, 1),
		unit_col:   lipgloss.NewStyle().Width(5).Padding(0, 1, 0, 0),
	}

	exportStyle := ExportStyle{
		colors:         colors,
		list:           lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForegroundBlend(primary, black, black, black, primary).BorderForegroundBlendOffset(10).BorderTop(true).BorderLeft(true),
		heading:        lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(highlight).Width(40),
		entry:          lipgloss.NewStyle().Padding(0, 1).Foreground(light_white).Width(40),
		selected_entry: lipgloss.NewStyle().Padding(0, 1).Foreground(light_white).Background(light_black).Width(40),
	}

	minWidth := 45
	minHeight := 5
	confirmMenuStyle := ConfirmMenuStyle{
		colors:        colors,
		minWidth:      minWidth - 2,
		minHeight:     minHeight,
		menu:          lipgloss.NewStyle().Width(minWidth).Height(minHeight).BorderStyle(lipgloss.NormalBorder()),
		description:   lipgloss.NewStyle().Width(minWidth - 2).Align(lipgloss.Center),
		confirmButton: lipgloss.NewStyle().Padding(0, 1).Margin(1, 0, 0, 6).Bold(true).Foreground(colors.black).Background(colors.primary),
		cancelButton:  lipgloss.NewStyle().Padding(0, 1).MarginTop(1).Bold(true).Foreground(colors.black).Background(colors.highlight),
	}

	return Style{
		colors: colors,

		docStyle: lipgloss.NewStyle().Padding(1, 2, 1, 2),
		title: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Italic(true).
			Foreground(black).
			Background(primary),

		listStyle:        listStyle,
		exportStyle:      exportStyle,
		confirmMenuStyle: confirmMenuStyle,
	}
}
