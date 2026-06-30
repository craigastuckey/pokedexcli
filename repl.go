package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/craigastuckey/pokedexcli/internal/location"
	"github.com/craigastuckey/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache) error
}

type config struct {
	next string
	prev string
}

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
	}

	return commands
}

func commandExit(conf *config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *pokecache.Cache) error {
	fmt.Println("Usage:\n\n\nhelp: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the map of the current location")
	fmt.Println("mapb: Displays the map of the previous location")
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache) error {
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

func commandMapb(conf *config, cache *pokecache.Cache) error {
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
