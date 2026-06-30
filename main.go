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

	commands := getCommands()

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
