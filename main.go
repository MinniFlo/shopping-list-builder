package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea	"charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if m, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		fmt.Printf("The last model state was: %v", m)
		os.Exit(1)
	}
}

func initialModel() model {
	cfg := loadConfig()
	mapping := loadMapping()
	recipes := BuildList(cfg, mapping)
	return model{
		collections:      recipes,
		collection_index: 0,
		entry_index:      -1,
		cfg:              cfg,
		mapping:          mapping,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.HandleDownMotion()
		case "k", "up":
			m.HandleUpMotion()
		case "J", "tab", "ctrl+d":
			if m.collection_index < len(m.collections)-1 {
				m.collection_index++
				m.entry_index = -1
			}
		case "K", "shift+tab", "ctrl+u":
			if m.collection_index > 0 {
				m.collection_index--
				m.entry_index = -1
			}
		case "+":
			m.CurrentRecipe().amount++
		case "-":
			if m.CurrentRecipe().amount > 1 {
				m.CurrentRecipe().amount--
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			number, _ := strconv.ParseInt(msg.String(), 10, 64)
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.category = CategoryFromInt(int(number))
			}
		case "l", "right":
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.state.Next()
			}
		case "h", "left":
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.state.Prev()
			}
		case "enter":
			renderShoppingListToFile(m.cfg.ShoppingListPath, m.ListEntriesForExport())
			m.UpdateMapping()
			saveMapping(m.mapping)
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	ri, ii := m.Indices()

	var s strings.Builder
	s.WriteString("Zutatenliste:\n\n")

	if len(m.collections) == 0 {
		s.WriteString("Es befinden sich keine Rezepte auf dem Essensplan …\n\n")
		fmt.Fprintf(&s, "Es befinden sich keine Rezepte auf dem Essensplan (%s) …\n\n", m.cfg.MealPlanPath)
	}

	for i, recipe := range m.collections {
		cursor_sym := " "
		if ri == i && ii == -1 {
			cursor_sym = ">"
		}

		fmt.Fprintf(&s, " %s  %d x %s\n", cursor_sym, recipe.amount, recipe.name)

		if i == ri {
			if len(recipe.entries) == 0 {
				s.WriteString("      Keine Zutaten gefunden …\n")
			}

			for j, incredient := range recipe.entries {
				if ii == j {
					cursor_sym = ">"
				} else {
					cursor_sym = " "
				}
				incredient_name := rightPadUnicodeConform(incredient.name, 30)
				incredient_amount := incredient.formated_amount(incredient.amount * float64(recipe.amount))
				category_symbol := incredient.category.Symbol()
				staged_string := ""
				switch incredient.state {
				case STAGED:
					staged_string = "[staged]"
				case MABY:
					staged_string = "[maby]"
				}

				fmt.Fprintf(&s, " %s    %s  %4s %-4s %s %s\n",
					cursor_sym,
					incredient_name,
					incredient_amount,
					incredient.unit,
					category_symbol,
					staged_string,
				)
			}
			s.WriteString("\n")
		}
	}
	s.WriteString("\n\nPress q to quit")

	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}
