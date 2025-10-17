package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	recipes          []recipe
	recipe_index     int
	incredient_index int
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

type recipe struct {
	name        string
	incredience []incredient
	amount      int
}

type category int

const (
	VEGI category = iota
	VEGTABLE
	COOL
	ASIA
	FROZEN
	PASTA
	MILK
	OTHER
	UNDEFINED
)

func (ca category) String() string {
	return [...]string{"Vegi-Regal", "Gemüse", "Kühl-Regal", "Asia-Regal", "TK-Regal", "Nudel-Regal", "Milch-Regal", "Gewürztes-Süßigkeiten-Regal", ""}[ca]
}

func (ca category) Symbol() string {
	return [...]string{"🌱", "🥕", "🧀", "🍙", "🧊", "🍝", "🥛", "🍫", "  "}[ca]
}

func categoryFromInt(i int) category {
	return category(i - 1)
}

type stage_state int

const (
	NOT_STAGED stage_state = iota
	MABY
	STAGED
)

const stage_state_max = 2

func (i *stage_state) Next() {
	if *i < stage_state_max {
		*i += 1
	}
}

func (i *stage_state) Prev() {
	if *i > 0 {
		*i -= 1
	}
}

type incredient struct {
	name     string
	amount   float32
	unit     string
	category category
	staged   stage_state
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
	for i, recipe := range m.recipes {
		cursor_sym := " "
		if ri == i && ii == -1 {
			cursor_sym = ">"
		}

		s += fmt.Sprintf(" %s  %d mal %s\n", cursor_sym, recipe.amount, recipe.name)

		if i == ri {
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
						" %s  | %s  %4s %-4s %s %s\n",
						cursor_sym,
						incredient_name,
						incredient_amount,
						incredient.unit,
						category_symbol,
						staged_string,
					)
			}
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

func buildIncredientData() []recipe {
	file, err := os.Open("resources/food/Essensplan.md")
	if err != nil {
		fmt.Printf("Failed to open 'Essensplan.md' with error: %v", err)
	}
	defer file.Close()

	var recipes []recipe
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		re := regexp.MustCompile(`\[\[(.*)\]\]`)
		match := re.FindStringSubmatch(row)

		if match != nil {
			recipe := recipe{name: match[1], amount: 1}
			recipes = append(recipes, recipe)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for i, recipe := range recipes {
		recipes[i].incredience = extractIncredientsFromRecipe(recipe.name)
	}

	return recipes
}

func extractIncredientsFromRecipe(recipe string) []incredient {
	path := fmt.Sprintf("resources/food/📝 Rezepte/%s.md", recipe)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open '%s' with error: %v", path, err)
	}
	defer file.Close()

	var incredience []incredient
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()

		re := regexp.MustCompile(`- \[.\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
		incredient_match := re.FindStringSubmatch(row)

		if incredient_match != nil {
			name := "INCREDIENT_MISSING"
			amount := 1.0
			unit := ""

			if value, err := strconv.ParseFloat(incredient_match[1], 32); err == nil {
				amount = value
			}

			if len(incredient_match[3]) > 0 {
				name = strings.TrimSpace(incredient_match[3])
			}

			if len(incredient_match[2]) > 0 {
				unit = strings.TrimSpace(incredient_match[2])
			}

			incredient := incredient{name: name, amount: float32(amount), unit: unit, category: UNDEFINED, staged: STAGED}
			incredience = append(incredience, incredient)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return incredience
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
