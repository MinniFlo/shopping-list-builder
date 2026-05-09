package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

func buildIncredientData(cfg config) []list_entry_collection {
	file, err := os.Open(cfg.MealPlanPath)
	if err != nil {
		fmt.Printf("Failed to open the Essensplan.md at '%v' with error: %v\n\n", cfg.MealPlanPath, err)
		os.Exit(1)
	}
	defer file.Close()

	var recipes []list_entry_collection
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		re := regexp.MustCompile(`- \[ \].*\[\[(.*)\]\]`)
		match := re.FindStringSubmatch(row)

		if match != nil {
			recipe := list_entry_collection{name: match[1], amount: 1}
			recipes = append(recipes, recipe)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for i, recipe := range recipes {
		recipes[i].entries = extractIncredientsFromRecipe(recipe.name, cfg.RecipesPath)
	}

	return recipes
}

func extractIncredientsFromRecipe(recipe string, base_path string) []list_entry {
	path := fmt.Sprintf("%s%s.md", base_path, recipe)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open the recipe with error: %v\n\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var incredience []list_entry
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

func createIncredientFromString(s string) (list_entry, error) {
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

		return list_entry{name: name, amount: float32(amount), unit: unit, category: UNDEFINED, staged: STAGED}, nil
	}

	return list_entry{}, errors.New("Invalid incredient string!")
}

func loadConfig() config {
	var cfg config

	var local_config = "./.shopping_list_builder.yml"
	local_data, err := os.ReadFile(local_config)
	if err == nil {
		if err := yaml.Unmarshal(local_data, &cfg); err == nil {
			return cfg
		} else {
			fmt.Printf("Found but could not parse local config file: %v\n\n", err)
		}
	}

	var user_config = "~/.config/shopping_list_builder.yml"
	user_data, err := os.ReadFile(user_config)
	if err == nil {
		if err := yaml.Unmarshal(user_data, &cfg); err == nil {
			return cfg
		} else {
			fmt.Printf("Found but could not parse user config file: %v\n\n", err)
		}
	}

	fmt.Printf("Could not find parsable config files at '%v' or '%v'", local_config, user_config)
	os.Exit(1)
	return cfg
}

func renderShoppingListToFile(shopping_list_file_path string, category_map map[category]map[string]*list_entry) {
	file, err := os.OpenFile(shopping_list_file_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var list_string strings.Builder
	list_string.WriteString("\n")

	for _, category := range categories {
		fmt.Fprintf(&list_string, "## %s\n", category.String())
		for _, entry := range category_map[category] {
			amount := strconv.FormatFloat(float64(entry.amount), 'f', -1, 64)
			fmt.Fprintf(&list_string, "- [ ] %s %s %s", amount, entry.unit, entry.name)
			if entry.staged == MABY {
				list_string.WriteString(" ?")
			}
			list_string.WriteString("\n")
		}
	}

	if _, err := file.WriteString(list_string.String()); err != nil {
		log.Fatal(err)
	}
}
