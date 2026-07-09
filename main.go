package main

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/craigastuckey/pokedexcli/internal/pokecache"
)

func main() {
	rl, err := readline.New("Pokedex > ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	conf := config{
		next: "https://pokeapi.co/api/v2/location-area/1/",
		prev: "",
	}

	commands := getCommands()
	cache := pokecache.NewCache(5)

	fmt.Println("Welcome to the Pokedex!")

	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF (Ctrl+D) or readline.ErrInterrupt (Ctrl+C)
			break
		}

		cleanInput := cleanInput(input)
		if len(cleanInput) == 0 {
			continue
		}

		if cmd, exists := commands[cleanInput[0]]; exists {
			args := make([]any, len(cleanInput[1:]))
			for i, arg := range cleanInput[1:] {
				args[i] = arg
			}
			cmd.callback(&conf, cache, args...)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
