package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

const mapping_path = "shopping_list_builder/mapping.yml"

var current_mapping_path = ""

func loadMapping() map[string]int {
	mapping := make(map[string]int)

	var local_config = fmt.Sprintf("./.%s", mapping_path)
	local_data, err := os.ReadFile(local_config)
	if err == nil {
		if err := yaml.Unmarshal(local_data, &mapping); err == nil {
			if mapping == nil {
				mapping = make(map[string]int)
			}
			current_mapping_path = local_config
			return mapping
		} else {
			fmt.Printf("Found but could not parse local mapping file: %v\n\n", err)
		}
	}

	var home_path, _ = os.UserHomeDir()
	var user_config = fmt.Sprintf("%s/.config/%s", home_path, mapping_path)
	user_data, err := os.ReadFile(user_config)
	if err == nil {
		if err := yaml.Unmarshal(user_data, &mapping); err == nil {
			current_mapping_path = user_config
			if mapping == nil {
				mapping = make(map[string]int)
			}
			return mapping
		} else {
			fmt.Printf("Found but could not parse user mapping file: %v\n\n", err)
		}
	}

	fmt.Printf("Could not find parsable mapping files at '%v' or '%v'", local_config, user_config)

	current_mapping_path = user_config
	return mapping
}

func saveMapping(mapping map[string]int) {
	file, err := os.OpenFile(current_mapping_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	out, err := yaml.Marshal(mapping)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(out)
	if err != nil {
		log.Fatal(err)
	}
}
