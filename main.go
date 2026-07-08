package main

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/craigastuckey/pokedexcli/internal/pokecache"
	"github.com/nexidian/gocliselect"
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
	menu := gocliselect.NewMenu("Continue or Exit?")
	menu.AddItem("Continue", "c")
	menu.AddItem("Exit", "e")

	choice := menu.Display()

	if choice == "e" {
		fmt.Println("Exiting...")
		return
	}

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
			cmd.callback(&conf, cache, cleanInput[1:]...)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
