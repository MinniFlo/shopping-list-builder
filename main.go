package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if m, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		fmt.Printf("The last model state was: %v", m)
		os.Exit(1)
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
			m.style = BuildStyle(m.bg_color, m.fg_color)
		}
	case tea.ForegroundColorMsg:
		m.fg_color = msg.String()
		if m.bg_color != "" {
			m.style = BuildStyle(m.bg_color, m.fg_color)
		}
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, default_key_map.Quit):
			return m, tea.Quit
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
		case key.Matches(msg, default_key_map.Confirm):
			renderShoppingListToFile(m.cfg.ShoppingListPath, m.ListEntriesForExport())
			m.UpdateMapping()
			saveMapping(m.mapping)
		case key.Matches(msg, default_key_map.Help):
			m.help.ShowAll = !m.help.ShowAll
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

	if len(m.collections) == 0 {
		doc.WriteString("Es befinden sich keine Rezepte auf dem Essensplan …\n\n")
		fmt.Fprintf(&doc, "Es befinden sich keine Rezepte auf dem Essensplan (%s) …\n\n", m.cfg.MealPlanPath)
	}

	ri, ii := m.Indices()

	for i, collection := range m.collections {
		header_is_selected := ri == i && ii == -1
		header := collection.buildHeader(header_is_selected, m.style)
		doc.WriteString(header)
		doc.WriteString("\n")

		if i == ri {
			if len(collection.entries) == 0 {
				doc.WriteString("Keine Zutaten gefunden …\n")
			}

			table := collection.buildTable(ii, m.style)
			doc.WriteString(m.style.table.Render(table))
			doc.WriteString("\n")
		}
	}

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
