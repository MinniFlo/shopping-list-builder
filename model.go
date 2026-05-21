package main

import (
	"errors"

	"charm.land/bubbles/v2/help"
)

type model struct {
	collections      []list_entry_collection
	collection_index int
	entry_index      int
	cfg              config
	mapping          map[string]int
	help             help.Model
}

func initialModel() model {
	cfg := loadConfig()
	mapping := loadMapping()
	recipes := BuildList(cfg, mapping)
	help := help.New()

	return model{
		collections:      recipes,
		collection_index: 0,
		entry_index:      -1,
		cfg:              cfg,
		mapping:          mapping,
		help:             help,
	}
}

func (m *model) ListEntriesForExport() map[category]map[string][]*list_entry {
	export_structure := make(map[category]map[string][]*list_entry)
	for _, category := range GetCategories() {
		export_structure[category] = make(map[string][]*list_entry)
	}

	for _, collection := range m.collections {
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

func (m *model) CurrentRecipe() *list_entry_collection {
	return &m.collections[m.collection_index]
}

func (m *model) CurrentIncredient() (*list_entry, error) {
	incredience := m.CurrentRecipe().entries
	if m.entry_index >= 0 && m.entry_index < len(incredience) {
		return &incredience[m.entry_index], nil
	}

	return nil, errors.New("No incredient selected")
}

func (m model) Indices() (int, int) {
	return m.collection_index, m.entry_index
}

func (m *model) HandleDownMotion() {
	ri, ii := m.Indices()

	switch {
	case ii < len(m.collections[ri].entries)-1:
		m.entry_index++
	case ii >= len(m.collections[ri].entries)-1 && ri < len(m.collections)-1:
		m.collection_index++
		m.entry_index = -1
	}
}

func (m *model) HandleUpMotion() {
	ri, ii := m.Indices()

	switch {
	case ii > -1:
		m.entry_index--
	case ii <= -1 && ri > 0:
		m.collection_index--
		m.entry_index = len(m.collections[m.collection_index].entries) - 1
	}
}

func (m *model) UpdateMapping() {
	for _, col := range m.collections {
		for _, entry := range col.entries {
			m.mapping[entry.name] = entry.category.ToInt()
		}
	}
}
