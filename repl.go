package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/craigastuckey/pokedexcli/internal/location"
	"github.com/craigastuckey/pokedexcli/internal/pokecache"
	"github.com/craigastuckey/pokedexcli/internal/pokemon"
	"github.com/nexidian/gocliselect"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, ...any) error
}

type config struct {
	next string
	prev string
}

var pokedex = make(map[string]pokemon.Pokemon)
var party = make([]pokemon.Pokemon, 6)

const locationAreaUrl string = "https://pokeapi.co/api/v2/location-area"
const pokemonUrl string = "https://pokeapi.co/api/v2/pokemon"

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
		"party": {
			name:        "party",
			description: "Show the list of pokemon in your party",
			callback:    commandParty,
		},
	}

	return commands
}

func commandExit(conf *config, cache *pokecache.Cache, args ...any) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *pokecache.Cache, args ...any) error {
	fmt.Println("Usage:\n\n\nhelp: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the map of the current location")
	fmt.Println("mapb: Displays the map of the previous location")
	fmt.Println("explore <location>: Explore the location passed in the argument")
	fmt.Println("catch <pokemon>: Catch the pokemon passed in the argument")
	fmt.Println("inspect <pokemon>: Inspect the pokemon passed in the argument")
	fmt.Println("pokedex: Displays the list of caught pokemon")
	fmt.Println("party: Displays the list of pokemon in your party")
	fmt.Println("party add <pokemon>: Add the pokemon passed in the argument to your party")
	fmt.Println("party remove <pokemon>: Remove the pokemon passed in the argument from your party")
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache, args ...any) error {
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
		conf.next = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID+1)
	}
	return nil
}

func commandMapb(conf *config, cache *pokecache.Cache, args ...any) error {
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
		temp.next = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID-20)
		conf.next = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID-20)
		if locationArea.ID-21 <= 0 {
			temp.prev = ""
			conf.prev = ""
		} else {
			temp.prev = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID-21)
			conf.prev = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID-21)
		}
		commandMap(&temp, cache)
	}
	return nil
}

func commandExplore(conf *config, cache *pokecache.Cache, args ...any) error {
	if len(args) == 0 {
		fmt.Println("Please provide a location to explore")
		return nil
	}

	var locationArea location.LocationArea
	var err error

	entry, exists := cache.Get(fmt.Sprintf("%s/%v/", locationAreaUrl, args[0]))
	if exists {
		locationArea = location.UnmarshalData(entry)
	} else {
		locationArea, err = location.GetLocationArea(fmt.Sprintf("%s/%v/", locationAreaUrl, args[0]))
		if err != nil {
			fmt.Println("Location not found")
			return fmt.Errorf("error fetching location: %w", err)
		}
		cache.Add(fmt.Sprintf("%s/%v/", locationAreaUrl, args[0]), location.MarshalData(locationArea))
	}

	menu := gocliselect.NewMenu(fmt.Sprintf("Exploring %s...\n", locationArea.Name))
	menu.AddItem("Explore next location", "next")
	menu.AddItem("Explore previous location", "prev")
	menu.AddItem("Encounter a Pokemon", "encounter")
	menu.AddItem("Back to map", "map")

	for {
		choice := menu.Display()
		switch choice {
		case "next":
			conf.prev = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID)
			conf.next = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID+2)
			next := locationArea.ID + 1
			commandExplore(conf, cache, next)
			return nil
		case "prev":
			conf.prev = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID-2)
			conf.next = fmt.Sprintf("%s/%d/", locationAreaUrl, locationArea.ID)
			prev := locationArea.ID - 1
			commandExplore(conf, cache, prev)
			return nil
		case "encounter":
			action := location.Encounter(locationArea)
			switch action[0] {
			case "throw":
				commandCatch(conf, cache, action[1])
			case "battle":
				fmt.Println("Battle feature not implemented yet")
			case "run":
				fmt.Println("You ran away from the wild", action[1])
			}
		case "map":
			return nil
		default:
			fmt.Println("Invalid choice")
			return fmt.Errorf("invalid explore choice")
		}
	}
}

func commandCatch(conf *config, cache *pokecache.Cache, args ...any) error {
	if len(args) == 0 {
		fmt.Println("Please provide a Pokemon to catch")
		return nil
	}

	var pm pokemon.Pokemon
	var err error

	entry, exists := cache.Get(fmt.Sprintf("%s/%v/", pokemonUrl, args[0]))
	if exists {
		pm = pokemon.UnmarshalData(entry)
	} else {
		pm, err = pokemon.GetPokemon(fmt.Sprintf("%s/%v/", pokemonUrl, args[0]))
		if err != nil {
			fmt.Println("Pokemon not found")
			return fmt.Errorf("error fetching Pokemon: %w", err)
		}
		cache.Add(fmt.Sprintf("%s/%v/", pokemonUrl, args[0]), pokemon.MarshalData(pm))
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pm.Name)

	if ok := pokemon.ThrowPokeball(pm); !ok {
		fmt.Printf("%s escaped!\n", pm.Name)
	} else {
		fmt.Printf("%s was caught!\n", pm.Name)
		fmt.Println("You may now inspect it with the inspect command")
		pokedex[pm.Name] = pm
		err = pokemon.AddToParty(&party, pm)
		if err != nil {
			fmt.Println("Party is full, pokemon added to your pokedex but not your party")
			return fmt.Errorf("party full: %w", err)
		}
	}

	return nil
}

func commandInspect(conf *config, cache *pokecache.Cache, args ...any) error {
	if len(args) == 0 {
		fmt.Println("Please provide a Pokemon to inspect")
		return nil
	}

	name, ok := args[0].(string)
	if !ok {
		fmt.Println("invalid pokemon name")
		return nil
	}

	pm, exists := pokedex[name]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	pokemon.GetStats(pm)
	return nil
}

func commandPokedex(conf *config, cache *pokecache.Cache, args ...any) error {
	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" -%s\n", name)
	}

	return nil
}

func commandParty(conf *config, cache *pokecache.Cache, args ...any) error {
	if len(args) == 0 {
		pokemon.GetParty(&party)
		return nil
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Please provide a Pokemon to add to your party")
			return nil
		}
		name, ok := args[1].(string)
		if !ok {
			fmt.Println("invalid pokemon name")
			return nil
		}

		pm, exists := pokedex[name]
		if !exists {
			fmt.Println("you have not caught that pokemon")
			return nil
		}

		err := pokemon.AddToParty(&party, pm)
		if err != nil {
			fmt.Println("Your party is full, remove a pokemon before adding another")
			return fmt.Errorf("error adding Pokemon to party: %w", err)
		}
		return nil
	case "remove":
		if len(args) < 2 {
			fmt.Println("Please provide a Pokemon to remove from your party")
			return nil
		}
		name, ok := args[1].(string)
		if !ok {
			fmt.Println("invalid pokemon name")
			return nil
		}
		pokemon.RemoveFromParty(&party, name)
		return nil
	default:
		fmt.Println("Unknown party command. Use 'party add <pokemon>' or 'party remove <pokemon>'")
	}
	return nil
}
