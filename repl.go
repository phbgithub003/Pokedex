package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous map",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
	}
}

func startRepl() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, ok := getCommands()[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		var err error
		if commandName == "explore" {
			if len(words) != 2 {
				fmt.Println("Explore command requires an area name")
				continue
			}
			err = command.callback(words[1])
		} else if commandName == "catch" {
			if len(words) != 2 {
				fmt.Println("Catch command requires a pokemon name")
				continue
			}
			err = command.callback(words[1])
		} else {
			err = command.callback()
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(...string) error {
	getLocationAreas()
	return nil
}

func commandMapb(...string) error {
	getPrevLocationArea()
	return nil
}

func commandExplore(areaName ...string) error {
	if len(areaName) == 0 {
		return fmt.Errorf("area name is required")
	}
	return getPokemonInArea(areaName[0])
}

func commandCatch(pokemonName ...string) error {
	if len(pokemonName) == 0 {
		return fmt.Errorf("pokemon name is required")
	}
	return catchPokemon(pokemonName[0])
}

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
}
