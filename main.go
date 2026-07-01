package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/craigastuckey/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	conf := config{
		next: "https://pokeapi.co/api/v2/location-area/1/",
		prev: "",
	}

	commands := getCommands()
	cache := pokecache.NewCache(5)

	fmt.Println("Welcome to the Pokedex!")

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cleanInput(input)

		if cmd, exists := commands[cleanInput[0]]; exists {
			cmd.callback(&conf, cache, cleanInput[1:]...)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
