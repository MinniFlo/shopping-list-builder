package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func buildIncredientData() []recipe {
	file, err := os.Open("resources/food/Essensplan.md")
	if err != nil {
		fmt.Printf("Failed to open the Essensplan.md with error: %v\n\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var recipes []recipe
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		re := regexp.MustCompile(`- \[ \].*\[\[(.*)\]\]`)
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
		fmt.Printf("Failed to open the recipe with error: %v\n\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var incredience []incredient
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()

		if inc, err := createIncredientFromString(row); err == nil {
			incredience = append(incredience, inc)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return incredience
}

func createIncredientFromString(s string) (incredient, error) {
	re := regexp.MustCompile(`- \[.\] ([0-9]+[.,][0-9]+|[0-9]+)?\s*(?i)(g|kg|l|ml|el|tl)?\b\s*(.*)`)
	incredient_match := re.FindStringSubmatch(s)

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

		return incredient{name: name, amount: float32(amount), unit: unit, category: UNDEFINED, staged: STAGED}, nil
	}

	return incredient{}, errors.New("Invalid incredient string!")
}

func createMealPlanContent() []section_content {
	var section_content []section_content
	return section_content
}
