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
)

func BuildList(cfg config, mapping map[string]int) []list_entry_collection {
	collections := extractRecipesFromMealPlan(cfg.MealPlanPath, cfg.RecipesPath, mapping)
	current_list_items, err := extractCurrentListEntriesFromShoppingList(cfg.ShoppingListPath, mapping)
	if err == nil {
		collections = append(collections, current_list_items)
	}

	return collections
}

func extractRecipesFromMealPlan(meal_plan_path string, recipes_path string, mapping map[string]int) []list_entry_collection {
	file, err := os.Open(meal_plan_path)
	if err != nil {
		fmt.Printf("Failed to open the Essensplan.md at '%v' with error: %v\n\n", meal_plan_path, err)
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
		path := fmt.Sprintf("%s%s.md", recipes_path, recipe.name)
		recipes[i].entries = extractListEntriesFromFile(path, mapping)
	}

	return recipes
}

func extractCurrentListEntriesFromShoppingList(shopping_list_path string, mapping map[string]int) (list_entry_collection, error) {
	entries := extractListEntriesFromFile(shopping_list_path, mapping)

	if len(entries) == 0 {
		return list_entry_collection{}, errors.New("No Items found on Soppinglist")
	}

	return list_entry_collection{name: "Zeugs von der Einkaufslist:", amount: 1, entries: entries}, nil
}

func extractListEntriesFromFile(path string, mapping map[string]int) []list_entry {
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

		if inc, err := createListEntryFromString(row, mapping); err == nil {
			incredience = append(incredience, inc)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return incredience
}

func renderShoppingListToFile(shopping_list_file_path string, category_map map[category]map[string]*list_entry) {
	if err := os.Truncate(shopping_list_file_path, 0); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(shopping_list_file_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var list_string strings.Builder

	for _, category := range GetCategories() {
		fmt.Fprintf(&list_string, "## %s\n", category.String())
		for _, entry := range category_map[category] {
			amount := strconv.FormatFloat(float64(entry.amount), 'f', -1, 64)
			fmt.Fprintf(&list_string, "- [ ] %s ", amount)
			if entry.unit != "" {
				fmt.Fprintf(&list_string, "%s ", entry.unit)
			}
			list_string.WriteString(entry.name)
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
