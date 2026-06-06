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
	style            Style
}

func initialListModel(collections []list_entry_collection) list_model {
	return list_model{
		collections:      collections,
		collection_index: 0,
		entry_index:      -1,
	}
}

func (m list_model) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, tea.RequestForegroundColor)
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
			m.CurrentRecipe().amount++
		case key.Matches(msg, default_key_map.Decrement):
			if m.CurrentRecipe().amount > 1 {
				m.CurrentRecipe().amount--
			}
		case key.Matches(msg, default_key_map.SetCategory):
			number, _ := strconv.ParseInt(msg.String(), 10, 64)
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.category = CategoryFromInt(int(number))
			}
		case key.Matches(msg, default_key_map.Right):
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.state.Next()
			}
		case key.Matches(msg, default_key_map.Left):
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.state.Prev()
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


func (m *list_model) CurrentRecipe() *list_entry_collection {
	return &m.collections[m.collection_index]
}

func (m *list_model) CurrentIncredient() (*list_entry, error) {
	incredience := m.CurrentRecipe().entries
	if m.entry_index >= 0 && m.entry_index < len(incredience) {
		return &incredience[m.entry_index], nil
	}

	return nil, errors.New("No incredient selected")
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
