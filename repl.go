package main

import (
	"fmt"
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
		"inspect": {
			name:        "inspect",
			description: "Inspect the pokemon passed in the argument",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show the list of caught pokemon",
			callback:    commandPokedex,
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
	fmt.Println("inspect <pokemon>: Inspect the pokemon passed in the argument")
	fmt.Println("pokedex: Displays the list of caught pokemon")
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache, args ...string) error {
	var locationArea location.LocationArea
	var err error

	for i := 0; i < 20; i++ {
		entry, exists := cache.Get(conf.next)
		if exists {
			locationArea = location.UnmarshalData(entry)
			fmt.Println(locationArea.Name)
		} else {
			locationArea, err = location.GetLocationArea(conf.next)
			if err != nil {
				fmt.Println("Location not found")
				return fmt.Errorf("error fetching location: %w", err)
			}
			cache.Add(conf.next, location.MarshalData(locationArea))
			fmt.Println(locationArea.Name)
		}

		conf.prev = conf.next
		conf.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID+1)
	}
	return nil
}

func commandMapb(conf *config, cache *pokecache.Cache, args ...string) error {
	locationArea, err := location.GetLocationArea(conf.next)
	if err != nil {
		fmt.Println("Location not found")
		return fmt.Errorf("error fetching location: %w", err)
	}

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
	var err error

	entry, exists := cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0]))
	if exists {
		locationArea = location.UnmarshalData(entry)
	} else {
		locationArea, err = location.GetLocationArea(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0]))
		if err != nil {
			fmt.Println("Location not found")
			return fmt.Errorf("error fetching location: %w", err)
		}
		cache.Add(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0]), location.MarshalData(locationArea))
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
	var err error

	entry, exists := cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]))
	if exists {
		pm = pokemon.UnmarshalData(entry)
	} else {
		pm, err = pokemon.GetPokemon(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]))
		if err != nil {
			fmt.Println("Pokemon not found")
			return fmt.Errorf("error fetching Pokemon: %w", err)
		}
		cache.Add(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0]), pokemon.MarshalData(pm))
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pm.Name)

	if ok := pokemon.ThrowPokeball(pm); !ok {
		fmt.Printf("%s escaped!\n", pm.Name)
	} else {
		fmt.Printf("%s was caught!\n", pm.Name)
		fmt.Println("You may now inspect it with the inspect command")
		pokedex[pm.Name] = pm
	}
	return nil
}

func commandInspect(conf *config, cache *pokecache.Cache, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a Pokemon to inspect")
		return nil
	}

	pm, exists := pokedex[args[0]]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	pokemon.GetStats(pm)
	return nil
}

func commandPokedex(conf *config, cache *pokecache.Cache, args ...string) error {
	fmt.Println("Your Pokedex:")
	for name, _ := range pokedex {
		fmt.Printf(" -%s\n", name)
	}

	return nil
}
