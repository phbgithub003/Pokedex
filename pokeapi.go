package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

// Create a global config variable
var pokeapiConfig Config

// GetConfig returns the global config
func GetConfig() *Config {
	return &pokeapiConfig
}

func getLocationAreas() error {
	url := "https://pokeapi.co/api/v2/location-area"
	if pokeapiConfig.Next != "" {
		url = pokeapiConfig.Next
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result LocationArea
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	pokeapiConfig.Next = result.Next
	if result.Previous != nil {
		pokeapiConfig.Previous = *result.Previous
	} else {
		pokeapiConfig.Previous = ""
	}

	for _, area := range result.Results {
		fmt.Printf("- %s\n", area.Name)
	}

	return nil
}

func getPrevLocationArea() error {
	if pokeapiConfig.Previous == "" {
		return fmt.Errorf("no previous location area")
	}

	url := pokeapiConfig.Previous
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result LocationArea
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	pokeapiConfig.Next = result.Next
	if result.Previous != nil {
		pokeapiConfig.Previous = *result.Previous
	} else {
		pokeapiConfig.Previous = ""
	}

	for _, area := range result.Results {
		fmt.Printf("- %s\n", area.Name)
	}

	return nil
}
