package main

import (
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
	export             export_model
}

func initialModel() model {
	cfg := loadConfig()
	mapping := loadMapping()
	collections := BuildList(cfg, mapping)
	help := help.New()
	list := initialListModel(collections)
	export := initialExportModel(collections, list.collection_index, list.entry_index)

	return model{
		cfg:      cfg,
		mapping:  mapping,
		help:     help,
		width:    0,
		height:   0,
		bg_color: "",
		fg_color: "",
		list:     list,
		export:   export,
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
			renderShoppingListToFile(m.cfg.ShoppingListPath, m.export.buildExportString(false, m.list.GetCurrentPositionID()))
			m.UpdateMapping()
			saveMapping(m.mapping)
		case key.Matches(msg, default_key_map.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			list, _ := m.list.Update(msg)
			m.list = list

			export, _ := m.export.Update(msg, m.list.collections, m.list.collection_index, m.list.entry_index)
			m.export = export
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
	}

	doc.WriteString(m.list.View().Content)

	help_view := m.style.docStyle.Render(m.help.View(default_key_map))

	join_view := lipgloss.JoinHorizontal(0, doc.String(), m.export.View(m.list.GetCurrentPositionID(), m.height).Content)

	document := m.style.docStyle.Render(join_view)

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
	m.list.style = style.listStyle
	m.export.style = style.exportStyle
}

func (m *model) UpdateMapping() {
	for _, col := range m.list.collections {
		for _, entry := range col.entries {
			m.mapping[entry.name] = entry.category.ToInt()
		}
	}
}
