package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type confirmation_model struct {
	description     string
	style           ConfirmMenuStyle
}

func initialConfrimationModel() confirmation_model {
	return confirmation_model{
		description: "Overwrite shopping list file with the generated list?",
	}
}

func (m confirmation_model) Init() tea.Cmd {
	return nil
}

func (m confirmation_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m confirmation_model) View() tea.View {
	description := m.style.description.Render(m.description)
	confirm := m.style.confirmButton.Render("Confirm")
	cancel := m.style.cancelButton.Render("Cancle")

	buttons := lipgloss.JoinHorizontal(0, cancel, confirm)
	centered_buttons := lipgloss.NewStyle().Width(m.style.minWidth).Align(lipgloss.Center).Render(buttons)
	menuContent := lipgloss.JoinVertical(0, description, centered_buttons)

	return tea.NewView(m.style.menu.Render(menuContent))
}
