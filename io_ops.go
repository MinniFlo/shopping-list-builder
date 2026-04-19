package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

func buildIncredientData(cfg config) []recipe {
	file, err := os.Open(cfg.MealPlanPath)
	if err != nil {
		fmt.Printf("Failed to open the Essensplan.md at '%v' with error: %v\n\n", cfg.MealPlanPath, err)
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
		recipes[i].incredience = extractIncredientsFromRecipe(recipe.name, cfg.RecipesPath)
	}

	return recipes
}

func extractIncredientsFromRecipe(recipe string, base_path string) []incredient {
	path := fmt.Sprintf("%s%s.md", base_path, recipe)

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
