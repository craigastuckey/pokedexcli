package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/craigastuckey/pokedexcli/internal/location"
	"github.com/craigastuckey/pokedexcli/internal/pokecache"
	"github.com/craigastuckey/pokedexcli/internal/pokemon"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, ...string) error
}

type config struct {
	next string
	prev string
}

var pokedex = make(map[string]pokemon.Pokemon)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the program",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the map of the current location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the map of the previous location",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the location passed in the argument",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch the pokemon passed in the argument",
			callback:    commandCatch,
		},
	}

	return commands
}

func commandExit(conf *config, cache *pokecache.Cache, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *pokecache.Cache, args ...string) error {
	fmt.Println("Usage:\n\n\nhelp: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the map of the current location")
	fmt.Println("mapb: Displays the map of the previous location")
	fmt.Println("explore <location>: Explore the location passed in the argument")
	fmt.Println("catch <pokemon>: Catch the pokemon passed in the argument")
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache, args ...string) error {
	var locationArea location.LocationArea

	for i := 0; i < 20; i++ {
		entry, exists := cache.Get(conf.next)
		if exists {
			locationArea = location.UnmarshalData(entry)
			fmt.Println(locationArea.Name)
		} else {
			locationArea = location.GetLocationArea(conf.next)
			cache.Add(conf.next, location.MarshalData(locationArea))
			fmt.Println(locationArea.Name)
		}

		conf.prev = conf.next
		conf.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID+1)
	}
	return nil
}

func commandMapb(conf *config, cache *pokecache.Cache, args ...string) error {
	locationArea := location.GetLocationArea(conf.next)

	if locationArea.ID <= 1 {
		fmt.Println("you're on the first page")
		return nil
	} else {
		temp := config{
			next: conf.next,
			prev: conf.prev,
		}
		temp.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID-20)
		conf.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID-20)
		if locationArea.ID-21 <= 0 {
			temp.prev = ""
			conf.prev = ""
		} else {
			temp.prev = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID-21)
			conf.prev = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID-21)
		}
		commandMap(&temp, cache)
	}
	return nil
}

func commandExplore(conf *config, cache *pokecache.Cache, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a location to explore")
		return nil
	}

	var locationArea location.LocationArea

	entry, exists := cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0]))
	if exists {
		locationArea = location.UnmarshalData(entry)
	} else {
		locationArea = location.GetLocationArea(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0]))
	}

	fmt.Printf("Exploring %s...\n", locationArea.Name)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(conf *config, cache *pokecache.Cache, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a Pokemon to catch")
		return nil
	}

	var pm pokemon.Pokemon

	entry, exists := cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]))
	if exists {
		pm = pokemon.UnmarshalData(entry)
	} else {
		pm = pokemon.GetPokemon(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]))
		cache.Add(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]), pokemon.MarshalData(pm))
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pm.Name)

	if ok := throwPokeball(pm); !ok {
		fmt.Printf("%s escaped!\n", pm.Name)
	} else {
		fmt.Printf("%s was caught!\n", pm.Name)
		pokedex[pm.Name] = pm
	}
	return nil
}

func throwPokeball(pm pokemon.Pokemon) bool {
	// Simulate a random chance of catching the Pokemon
	catchChance := 0.5 // 50% chance to catch
	if pm.BaseExperience > 100 {
		catchChance = 0.3 // Harder to catch if base experience is high
	}

	return rand.Float64() < catchChance
}
