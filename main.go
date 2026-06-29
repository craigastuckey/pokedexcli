package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	conf := config{
		next: "https://pokeapi.co/api/v2/location-area/1/",
		prev: "",
	}

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

	fmt.Println("Welcome to the Pokedex!")

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cleanInput(input)

		if cmd, exists := commands[cleanInput[0]]; exists {
			cmd.callback(&conf)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
