package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	cfg                config
	mapping            map[string]int
	help               help.Model
	width, height      int
	bg_color, fg_color string
	style              Style
	list               list_model
}

func initialModel() model {
	cfg := loadConfig()
	mapping := loadMapping()
	collections := BuildList(cfg, mapping)
	help := help.New()
	list := list_model{
		collections:      collections,
		collection_index: 0,
		entry_index:      -1,
	}

	return model{
		cfg:              cfg,
		mapping:          mapping,
		help:             help,
		width:            0,
		height:           0,
		bg_color:         "",
		fg_color:         "",
		list:             list,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, tea.RequestForegroundColor)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.BackgroundColorMsg:
		m.bg_color = msg.String()
		if m.fg_color != "" {
			m.setStyle()
		}
	case tea.ForegroundColorMsg:
		m.fg_color = msg.String()
		if m.bg_color != "" {
			m.setStyle()
		}
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, default_key_map.Quit):
			return m, tea.Quit
		case key.Matches(msg, default_key_map.Confirm):
			renderShoppingListToFile(m.cfg.ShoppingListPath, ListEntriesForExport(m.list.collections))
			m.UpdateMapping()
			saveMapping(m.mapping)
		case key.Matches(msg, default_key_map.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			list, _ := m.list.Update(msg)
			m.list = list
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	var doc strings.Builder

	if m.bg_color == "" || m.fg_color == "" {
		doc.WriteString("Loading ...")
		return tea.NewView(doc.String())
	}

	doc.WriteString(m.style.title.Render("SHOPPING LIST BUILDER"))
	doc.WriteString("\n\n")

	if len(m.list.collections) == 0 {
		doc.WriteString("Es befinden sich keine Rezepte auf dem Essensplan …\n\n")
		fmt.Fprintf(&doc, "Es befinden sich keine Rezepte auf dem Essensplan (%s) …\n\n", m.cfg.MealPlanPath)
	}

	doc.WriteString(m.list.View().Content)

	help_view := m.style.docStyle.Render(m.help.View(default_key_map))
	document := m.style.docStyle.Render(doc.String())

	help_height := strings.Count(help_view, "\n")
	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(document),
		lipgloss.NewLayer(help_view).Y(m.height - (help_height)),
	}

	comp := lipgloss.NewCompositor(layers...)

	v := tea.NewView(lipgloss.Sprintln(comp.Render()))
	v.AltScreen = true
	return v
}

func (m *model) setStyle() {
	style := BuildStyle(m.bg_color, m.fg_color)
	m.style = style
	m.list.style = style
}

func ListEntriesForExport(collections []list_entry_collection) map[category]map[string][]*list_entry {
	export_structure := make(map[category]map[string][]*list_entry)
	for _, category := range GetCategories() {
		export_structure[category] = make(map[string][]*list_entry)
	}

	for _, collection := range collections {
		for _, entry := range collection.entries {
			if entry.state == NOT_STAGED {
				continue
			}

			entry.amount *= float64(collection.amount)

			export_entry_list, ok := export_structure[entry.category][entry.name]
			if ok == false {
				export_structure[entry.category][entry.name] = []*list_entry{
					{
						name:     entry.name,
						unit:     entry.unit,
						amount:   entry.amount,
						category: entry.category,
						state:    entry.state,
					},
				}
			} else {
				export_structure[entry.category][entry.name] = mergeOrInsertIntoExportEntryList(export_entry_list, entry)
			}
		}
	}

	return export_structure
}

func mergeOrInsertIntoExportEntryList(export_entry_list []*list_entry, entry list_entry) []*list_entry {
	var transformation_info transformation_info
	var exact_match bool = false
	var target_entry *list_entry
	for _, export_entry := range export_entry_list {
		if export_entry.unit == entry.unit {
			exact_match = true
			target_entry = export_entry
			break
		}
		ti, ok := unit_transformations[MakeUnitPair(entry.unit, export_entry.unit)]
		if ok {
			target_entry = export_entry
			transformation_info = ti
		}
	}
	if target_entry != nil {
		if exact_match {
			target_entry.amount += entry.amount
		} else {
			unit_values := map[unit]float64{target_entry.unit: target_entry.amount, entry.unit: entry.amount}
			new_unit, new_amount := MergeUnitAmounts(unit_values, transformation_info)

			target_entry.unit = new_unit
			target_entry.amount = new_amount
		}

		target_entry.state = MergeStagingState(target_entry.state, entry.state)
	} else {
		export_entry_list = append(
			export_entry_list,
			&list_entry{
				name:     entry.name,
				unit:     entry.unit,
				amount:   entry.amount,
				category: entry.category,
				state:    entry.state,
			},
		)
	}

	return export_entry_list
}

func (m *model) UpdateMapping() {
	for _, col := range m.list.collections {
		for _, entry := range col.entries {
			m.mapping[entry.name] = entry.category.ToInt()
		}
	}
}
