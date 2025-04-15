package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type LocationAreaPokemon struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Create a global config variable
var pokeapiConfig Config

// Create a global cache with a 5-minute expiration
var pokeapiCache = NewCache(5 * time.Minute)

// GetConfig returns the global config
func GetConfig() *Config {
	return &pokeapiConfig
}

// fetchFromCacheOrRemote gets data from cache if available or makes a remote request
func fetchFromCacheOrRemote(url string) ([]byte, error) {
	// Check if we have this URL cached
	if cachedData, ok := pokeapiCache.Get(url); ok {
		return cachedData, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Cache the response
	pokeapiCache.Add(url, body)
	return body, nil
}

// processLocationAreaData parses the API response and updates config
func processLocationAreaData(body []byte) error {
	var result LocationArea
	err := json.Unmarshal(body, &result)
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

func getLocationAreas() error {
	url := "https://pokeapi.co/api/v2/location-area"
	if pokeapiConfig.Next != "" {
		url = pokeapiConfig.Next
	}

	body, err := fetchFromCacheOrRemote(url)
	if err != nil {
		return err
	}

	return processLocationAreaData(body)
}

func getPrevLocationArea() error {
	if pokeapiConfig.Previous == "" {
		return fmt.Errorf("you're on the first page")
	}

	body, err := fetchFromCacheOrRemote(pokeapiConfig.Previous)
	if err != nil {
		return err
	}

	return processLocationAreaData(body)
}

func getPokemonInArea(areaName string) error {
	fmt.Printf("Exploring %s...\n", areaName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	body, err := fetchFromCacheOrRemote(url)
	if err != nil {
		return err
	}

	var result LocationAreaPokemon
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range result.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}
