package main

import (
	"errors"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type list_model struct {
	collections      []list_entry_collection
	collection_index int
	entry_index      int
	style            ListStyle
}

func initialListModel(collections []list_entry_collection) list_model {
	return list_model{
		collections:      collections,
		collection_index: 0,
		entry_index:      -1,
	}
}

func (m list_model) Init() tea.Cmd {
	return nil
}

func (m list_model) Update(msg tea.Msg) (list_model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, default_key_map.Down):
			m.HandleDownMotion()
		case key.Matches(msg, default_key_map.Up):
			m.HandleUpMotion()
		case key.Matches(msg, default_key_map.SectionDown):
			if m.collection_index < len(m.collections)-1 {
				m.collection_index++
				m.entry_index = -1
			}
		case key.Matches(msg, default_key_map.SectionUp):
			if m.collection_index > 0 {
				m.collection_index--
				m.entry_index = -1
			}
		case key.Matches(msg, default_key_map.Increment):
			m.currentCollection().amount++
		case key.Matches(msg, default_key_map.Decrement):
			if m.currentCollection().amount > 1 {
				m.currentCollection().amount--
			}
		case key.Matches(msg, default_key_map.SetCategory):
			number, _ := strconv.ParseInt(msg.String(), 10, 64)
			if entry, err := m.currentEntry(); err == nil {
				entry.category = CategoryFromInt(int(number))
			} else {
				m.batchSetCategory(int(number))
			}
		case key.Matches(msg, default_key_map.Right):
			if inc, err := m.currentEntry(); err == nil {
				inc.state.Next()
			} else {
				m.batchIncrementState()
			}
		case key.Matches(msg, default_key_map.Left):
			if inc, err := m.currentEntry(); err == nil {
				inc.state.Prev()
			} else {
				m.batchDecrementState()
			}
		}
	}

	return m, nil
}

func (m list_model) View() tea.View {
	var list strings.Builder

	ri, ii := m.Indices()

	for i, collection := range m.collections {
		header_is_selected := ri == i && ii == -1
		header := collection.buildHeader(header_is_selected, m.style)
		list.WriteString(header)
		list.WriteString("\n")

		if i == ri {
			if len(collection.entries) == 0 {
				list.WriteString("Keine Zutaten gefunden …\n")
			}

			table := collection.buildTable(ii, m.style)
			list.WriteString(m.style.table.Render(table))
			list.WriteString("\n")
		}
	}

	v := tea.NewView(list.String())
	return v
}

func (m *list_model) currentCollection() *list_entry_collection {
	return &m.collections[m.collection_index]
}

func (m *list_model) currentEntry() (*list_entry, error) {
	entry := m.currentCollection().entries
	if m.entry_index >= 0 && m.entry_index < len(entry) {
		return &entry[m.entry_index], nil
	}

	return nil, errors.New("No entry selected")
}

func (m list_model) Indices() (int, int) {
	return m.collection_index, m.entry_index
}

func (m *list_model) HandleDownMotion() {
	ri, ii := m.Indices()

	switch {
	case ii < len(m.collections[ri].entries)-1:
		m.entry_index++
	case ii >= len(m.collections[ri].entries)-1 && ri < len(m.collections)-1:
		m.collection_index++
		m.entry_index = -1
	}
}

func (m *list_model) HandleUpMotion() {
	ri, ii := m.Indices()

	switch {
	case ii > -1:
		m.entry_index--
	case ii <= -1 && ri > 0:
		m.collection_index--
		m.entry_index = len(m.collections[m.collection_index].entries) - 1
	}
}

func (m *list_model) batchSetCategory(category_number int) {
	coll := m.currentCollection()
	for i := range coll.entries {
		coll.entries[i].category = CategoryFromInt(category_number)
	}
}

func (m *list_model) batchIncrementState() {
	coll := m.currentCollection()
	for i := range coll.entries {
		coll.entries[i].state.Next()
	}
}

func (m *list_model) batchDecrementState() {
	coll := m.currentCollection()
	for i := range coll.entries {
		coll.entries[i].state.Prev()
	}
}
