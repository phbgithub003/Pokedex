package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

type PokemonStrength struct {
	BaseEexperience int `json:base_experience`
}

// Updated Pokemon struct with additional fields
type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Create a global config variable
var pokeapiConfig Config

// Create a global cache with a 5-minute expiration
var pokeapiCache = NewCache(5 * time.Minute)

// Create a map to store caught Pokemon
var pokedex = make(map[string]Pokemon)

// Add a initialization function for the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// Add a function to inspect caught Pokemon
func inspectPokemon(pokemonName string) error {
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("you haven't caught %s yet", pokemonName)
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  -%s\n", typeInfo.Type.Name)
	}

	return nil
}

// Add a function to list all caught Pokemon
func listCaughtPokemon() error {
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}

// Update the catchPokemon function to store caught Pokemon
func catchPokemon(pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	body, err := fetchFromCacheOrRemote(url)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	// Calculate catch probability based on base_experience
	// Higher base_experience = harder to catch
	catchRate := 100 - pokemon.BaseExperience/2
	if catchRate < 10 {
		catchRate = 10 // Minimum 10% chance
	} else if catchRate > 90 {
		catchRate = 90 // Maximum 90% chance
	}

	// Generate a random number between 1 and 100
	randomNum := rand.Intn(100) + 1

	fmt.Printf("Trying to catch with %d%% chance...\n", catchRate)

	// Compare the random number with the catch rate
	if randomNum <= catchRate {
		fmt.Printf("Congratulations! You caught %s!\n", pokemon.Name)
		// Add the Pokemon to the Pokedex
		pokedex[pokemon.Name] = pokemon
		return nil
	} else {
		fmt.Printf("Oh no! %s escaped!\n", pokemon.Name)
		return nil
	}
}
