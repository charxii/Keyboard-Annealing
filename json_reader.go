package main

import (
	"encoding/json"
	"os"
)

type LayoutMap map[string]string

func saveLayoutToJSON(filename string, layouts LayoutMap) error {
	data, err := json.MarshalIndent(layouts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// layouts are in row form, but will be converted to column form
func loadLayoutFromJSON(filename string) (LayoutMap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var layouts LayoutMap
	err = json.Unmarshal(data, &layouts)
	if err != nil {
		return nil, err
	}

	return layouts, nil
}

type StatMap map[string]map[string]float64

func saveStatsToJSON(filename string, stats StatMap) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
