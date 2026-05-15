package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type path_config struct {
	MealPlanPath     string `yaml:"meal_plan_path"`
	RecipesPath      string `yaml:"recipes_path"`
	ShoppingListPath string `yaml:"shopping_list_path"`
}

func loadConfig() path_config {
	var cfg path_config

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
