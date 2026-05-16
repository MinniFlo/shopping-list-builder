package main

type unit int

const (
	None unit = iota
	Gram
	Kilogram
	Milliliter
	Liter
	Teeloeffel
	Essloeffel
)

var unit_disply_strings = map[unit]string{
	None:       "",
	Gram:       "g",
	Kilogram:   "kg",
	Milliliter: "ml",
	Liter:      "l",
	Teeloeffel: "Tl",
	Essloeffel: "El",
}

var string_unit_mapping = map[string]unit{
	"g":  Gram,
	"kg": Kilogram,
	"ml": Milliliter,
	"l":  Liter,
	"tl": Teeloeffel,
	"el": Essloeffel,
}

func (u unit) String() string {
	return unit_disply_strings[u]
}

func UnitFromString(s string) unit {
	if unit, ok := string_unit_mapping[s]; ok {
		return unit
	}

	return None
}

type unit_pair struct{ u1, u2 unit }

type transformation_info struct {
	unit_to_transform unit
	target_unit       unit
	scale_factor      float64
}

func MakeUnitPair(u1, u2 unit) unit_pair {
	if u1 <= u2 {
		return unit_pair{u1: u1, u2: u2}
	}
	return unit_pair{u1: u2, u2: u1}
}

var unit_transformations = map[unit_pair]transformation_info{
	{Gram, Kilogram}: {
		unit_to_transform: Gram, target_unit: Kilogram, scale_factor: 0.001,
	},
	{Milliliter, Liter}: {
		unit_to_transform: Milliliter, target_unit: Liter, scale_factor: 0.001,
	},
	{Teeloeffel, Essloeffel}: {
		unit_to_transform: Teeloeffel, target_unit: Essloeffel, scale_factor: 0.5,
	},
}

func MergeUnitAmounts(unit_values map[unit]float64, ti transformation_info) (unit, float64) {
	transform_unit_value := unit_values[ti.unit_to_transform]
	target_unit_value := unit_values[ti.target_unit]
	new_value := RoundToThreeDigitsAfterPeriode(transform_unit_value*ti.scale_factor) + target_unit_value

	return ti.target_unit, new_value
}
