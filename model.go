package main

import "errors"

type model struct {
	collections      []list_entry_collection
	collection_index int
	entry_index      int
	cfg              path_config
}

func (m *model) ListEntriesForExport() map[category]map[string]*list_entry {
	export_structure := make(map[category]map[string]*list_entry)
	for _, category := range GetCategories() {
		export_structure[category] = make(map[string]*list_entry)
	}

	for _, collection := range m.collections {
		for _, entry := range collection.entries {
			if entry.staged == NOT_STAGED {
				continue
			}
			if export_entry, ok := export_structure[entry.category][entry.name]; ok == false {
				export_structure[entry.category][entry.name] =
					&list_entry{
						name:     entry.name,
						unit:     entry.unit,
						amount:   entry.amount * float32(collection.amount),
						category: entry.category,
						staged:   entry.staged,
					}
			} else {
				// TODO: unify units
				export_entry.amount += entry.amount * float32(collection.amount)
			}
		}
	}

	return export_structure
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
