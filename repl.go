package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Usage:\n\n\nhelp: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the map of the current location")
	fmt.Println("mapb: Displays the map of the previous location")
	return nil
}

func commandMap(conf *config) error {
	for i := 0; i < 20; i++ {
		res, err := http.Get(conf.next)
		if err != nil {
			fmt.Println("Error fetching data:", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		if res.StatusCode != http.StatusOK {
			fmt.Println("Error: received non-OK HTTP status:", res.Status)
		}
		res.Body.Close()

		var locationArea LocationArea
		err = json.Unmarshal(body, &locationArea)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		}

		fmt.Println(locationArea.Name)

		conf.prev = conf.next
		conf.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationArea.ID+1)
	}
	return nil
}

func commandMapb(conf *config) error {
	res, err := http.Get(conf.next)
	if err != nil {
		fmt.Println("Error fetching data:", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-OK HTTP status:", res.Status)
	}
	res.Body.Close()

	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
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
		commandMap(&temp)
	}
	return nil
}
