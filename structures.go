package main

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

var categories = []category{VEGI, VEGTABLE, COOL, ASIA, FROZEN, PASTA, MILK, OTHER, UNDEFINED}

func (ca category) String() string {
	return [...]string{"Vegi-Regal", "Gemüse", "Kühl-Regal", "Asia-Regal", "TK-Regal", "Nudel-Regal", "Milch-Regal", "Gewürztes-Süßigkeiten-Regal", "Stuff"}[ca]
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

type list_entry struct {
	name     string
	amount   float32
	unit     string
	category category
	staged   stage_state
}

type list_entry_collection struct {
	name    string
	entries []list_entry
	amount  int
}

type config struct {
	MealPlanPath     string `yaml:"meal_plan_path"`
	RecipesPath      string `yaml:"recipes_path"`
	ShoppingListPath string `yaml:"shopping_list_path"`
}
