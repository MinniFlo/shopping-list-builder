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

type recipe struct {
	name        string
	incredience []incredient
	amount      int
}

type section_content struct {
	name         string
	catecategory category
	content      string
}
