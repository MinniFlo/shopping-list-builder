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

func GetCategories() []category {
	return []category{VEGI, VEGTABLE, COOL, ASIA, FROZEN, PASTA, MILK, OTHER, UNDEFINED}
}

func (ca category) String() string {
	return [...]string{"Vegi-Regal", "Gemüse", "Kühl-Regal", "Asia-Regal", "TK-Regal", "Nudel-Regal", "Milch-Regal", "Gewürztes-Süßigkeiten-Regal", "Stuff"}[ca]
}

func (ca category) Symbol() string {
	return [...]string{"🌱", "🥕", "🧀", "🍙", "🧊", "🍝", "🥛", "🍫", "  "}[ca]
}

func CategoryFromInt(i int) category {
	return category(i - 1)
}
