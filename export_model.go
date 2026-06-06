package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type export_model struct {
	structure map[category][]*export_entry
	style     ExportStyle
}

func initialExportModel(collections []list_entry_collection) export_model {

	structure := ListEntriesForExport(collections)

	return export_model{structure: structure}
}

func (m export_model) Init() tea.Cmd {
	return nil
}

func (m export_model) Update(msg tea.Msg) (export_model, tea.Cmd) {
	// switch msg := msg.(type) {}

	return m, nil
}

func (m export_model) View(selected_entry_id int) tea.View {
	var preview strings.Builder

	list := m.style.list.Render(m.buildExportString(true, selected_entry_id))
	preview.WriteString(list)

	v := tea.NewView(preview.String())
	return v
}

func ListEntriesForExport(collections []list_entry_collection) map[category][]*export_entry {
	work_structure := make(map[category]map[string][]*export_entry)
	for _, category := range GetCategories() {
		work_structure[category] = make(map[string][]*export_entry)
	}

	for _, collection := range collections {
		for _, entry := range collection.entries {
			if entry.state == NOT_STAGED {
				continue
			}

			entry.amount *= float64(collection.amount)

			export_entry_list, ok := work_structure[entry.category][entry.name]
			if ok == false {
				work_structure[entry.category][entry.name] = []*export_entry{
					{
						name:           entry.name,
						unit:           entry.unit,
						amount:         entry.amount,
						state:          entry.state,
						list_entry_ids: []int{entry.id},
					},
				}
			} else {
				work_structure[entry.category][entry.name] = mergeOrInsertIntoExportEntryList(export_entry_list, entry)
			}
		}
	}

	return flattenExportStrucutre(work_structure)
}

func mergeOrInsertIntoExportEntryList(export_entry_list []*export_entry, entry list_entry) []*export_entry {
	var transformation_info transformation_info
	var exact_match bool = false
	var target_entry *export_entry
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
		target_entry.list_entry_ids = append(target_entry.list_entry_ids, entry.id)
		slices.Sort(target_entry.list_entry_ids)
	} else {
		export_entry_list = append(
			export_entry_list,
			&export_entry{
				name:           entry.name,
				unit:           entry.unit,
				amount:         entry.amount,
				state:          entry.state,
				list_entry_ids: []int{entry.id},
			},
		)
	}
	return export_entry_list
}

func flattenExportStrucutre(structure map[category]map[string][]*export_entry) map[category][]*export_entry {
	export_structure := make(map[category][]*export_entry)

	for category, name_map := range structure {
		export_structure[category] = []*export_entry{}
		for _, export_entries := range name_map {
			export_structure[category] = append(export_structure[category], export_entries...)
		}
		sort.Sort(ByListEntryIds(export_structure[category]))
	}

	return export_structure
}

func (e export_model) buildExportString(with_styles bool, selected_entry_id int) string {
	heading_style := lipgloss.NewStyle()
	entry_style := lipgloss.NewStyle()
	selected_entry_style := lipgloss.NewStyle()

	if with_styles {
		heading_style = e.style.heading
		entry_style = e.style.entry
		selected_entry_style = e.style.selected_entry
	}

	var list_string strings.Builder

	for _, category := range GetCategories() {
		heading := heading_style.Render(headingString(category.String()))
		list_string.WriteString(heading)
		list_string.WriteString("\n")
		for _, entry := range e.structure[category] {
			selected := false
			if with_styles {
				for _, id := range entry.list_entry_ids {
					if id == selected_entry_id {
						selected = true
					}
				}
			}
			list_entry := checkboxListEntryString(entry)
			if selected {
				list_string.WriteString(selected_entry_style.Render(list_entry))
			} else {
				list_string.WriteString(entry_style.Render(list_entry))
			}
			list_string.WriteString("\n")

		}
	}

	return list_string.String()
}

func checkboxListEntryString(entry *export_entry) string {
	var entry_string strings.Builder

	entry_string.WriteString("- [ ] ")
	if amount := formatedAmount(entry.amount); amount != "" {
		fmt.Fprintf(&entry_string, "%s ", amount)
	}
	if entry.unit != None {
		fmt.Fprintf(&entry_string, "%s ", entry.unit)
	}
	entry_string.WriteString(entry.name)
	if entry.state == MABY {
		entry_string.WriteString(" ?")
	}

	return entry_string.String()
}

func headingString(heading string) string {
	return fmt.Sprintf("### %s", heading)
}
