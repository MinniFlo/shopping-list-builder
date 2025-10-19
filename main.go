package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	recipes           []recipe
	recipe_index      int
	incredient_index  int
	meal_plan_content []section_content
}

func (m *model) CurrentRecipe() *recipe {
	return &m.recipes[m.recipe_index]
}

func (m *model) CurrentIncredient() (*incredient, error) {
	incredience := m.CurrentRecipe().incredience
	if m.incredient_index >= 0 && m.incredient_index < len(incredience) {
		return &incredience[m.incredient_index], nil
	}

	return nil, errors.New("No incredient selected")
}

func (m model) Indices() (int, int) {
	return m.recipe_index, m.incredient_index
}

func initialModel() model {
	recipes := buildIncredientData()
	return model{
		recipes:          recipes,
		recipe_index:     0,
		incredient_index: -1,
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
			if m.recipe_index < len(m.recipes)-1 {
				m.recipe_index++
				m.incredient_index = -1
			}
		case "K", "shift+tab", "ctrl+u":
			if m.recipe_index > 0 {
				m.recipe_index--
				m.incredient_index = -1
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
				inc.category = categoryFromInt(int(number))
			}
		case "l", "right":
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.staged.Next()
			}
		case "h", "left":
			if inc, err := m.CurrentIncredient(); err == nil {
				inc.staged.Prev()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	ri, ii := m.Indices()

	s := "Zutatenliste:\n\n"

	if len(m.recipes) == 0 {
		s += "Es befinden sich keine Rezepte auf dem Essensplan …\n\n"
	}

	for i, recipe := range m.recipes {
		cursor_sym := " "
		if ri == i && ii == -1 {
			cursor_sym = ">"
		}

		s += fmt.Sprintf(" %s  %d x %s\n", cursor_sym, recipe.amount, recipe.name)

		if i == ri {
			if len(recipe.incredience) == 0 {
				s += "      Keine Zutaten gefunden …\n"
			}

			for j, incredient := range recipe.incredience {
				if ii == j {
					cursor_sym = ">"
				} else {
					cursor_sym = " "
				}
				incredient_name := rightPadUnicodeConform(incredient.name, 30)
				incredient_amount := strconv.FormatFloat(float64(incredient.amount*float32(recipe.amount)), 'f', -1, 64)
				category_symbol := incredient.category.Symbol()
				staged_string := ""
				switch incredient.staged {
				case STAGED:
					staged_string = "[staged]"
				case MABY:
					staged_string = "[maby]"
				}

				s +=
					fmt.Sprintf(
						" %s    %s  %4s %-4s %s %s\n",
						cursor_sym,
						incredient_name,
						incredient_amount,
						incredient.unit,
						category_symbol,
						staged_string,
					)
			}
			s += "\n"
		}
	}
	s += "\n\nPress q to quit"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if m, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		fmt.Printf("The last model state was: %v", m)
		os.Exit(1)
	}
}

func (m *model) HandleDownMotion() {
	ri, ii := m.Indices()

	switch {
	case ii < len(m.recipes[ri].incredience)-1:
		m.incredient_index++
	case ii >= len(m.recipes[ri].incredience)-1 && ri < len(m.recipes)-1:
		m.recipe_index++
		m.incredient_index = -1
	}
}

func (m *model) HandleUpMotion() {
	ri, ii := m.Indices()

	switch {
	case ii > -1:
		m.incredient_index--
	case ii <= -1 && ri > 0:
		m.recipe_index--
		m.incredient_index = len(m.recipes[m.recipe_index].incredience) - 1
	}
}

func rightPadUnicodeConform(s string, pad_value int) string {
	pad_amt := pad_value - utf8.RuneCountInString(s)

	switch {
	case pad_amt == 0:
		return s
	case pad_amt > 0:
		runes := []rune(s)
		for range pad_amt {
			runes = append(runes, rune(' '))
		}
		return string(runes)
	case pad_amt < 0:
		runes := []rune(s)
		runes = append(runes[:pad_value-1], rune('…'))
		return string(runes)
	}

	return s
}
