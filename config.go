package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type config struct {
	MealPlanPath     string `yaml:"meal_plan_path"`
	RecipesPath      string `yaml:"recipes_path"`
	ShoppingListPath string `yaml:"shopping_list_path"`
}

const config_path = "shopping_list_builder/config.yml"

func loadConfig() config {
	var cfg config

	var local_config = fmt.Sprintf("./.%s", config_path)
	local_data, err := os.ReadFile(local_config)
	if err == nil {
		if err := yaml.Unmarshal(local_data, &cfg); err == nil {
			return cfg
		} else {
			fmt.Printf("Found but could not parse local config file: %v\n\n", err)
		}
	}

	var home_path, _ = os.UserHomeDir()
	var user_config = fmt.Sprintf("%s/.config/%s", home_path, config_path)
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
