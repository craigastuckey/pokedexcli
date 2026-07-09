package location

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/nexidian/gocliselect"
)

type LocationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationArea(url string) (LocationArea, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error fetching data: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error reading response body: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("error: received non-OK HTTP status: %s", res.Status)
	}
	res.Body.Close()

	locationArea := UnmarshalData(body)

	return locationArea, nil
}

func UnmarshalData(data []byte) LocationArea {
	var locationArea LocationArea
	err := json.Unmarshal(data, &locationArea)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	return locationArea
}

func MarshalData(locationArea LocationArea) []byte {
	data, err := json.Marshal(locationArea)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	return data
}

func Encounter(locationArea LocationArea) {
	fmt.Printf("Encountering Pokemon")
	ch := make(chan struct{})
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
		fmt.Println()
		ch <- struct{}{}
	}()

	pokemon := locationArea.PokemonEncounters[rand.Intn(len(locationArea.PokemonEncounters))].Pokemon
	menu := gocliselect.NewMenu(fmt.Sprintf("You encountered a wild %s!\n", pokemon.Name))
	menu.AddItem("Throw a Pokeball", "throw")
	menu.AddItem("Battle", "battle")
	menu.AddItem("Run away", "run")

	<-ch

	choice := menu.Display()
	switch choice {
	case "throw":
		fmt.Printf("You threw a Pokeball at %s!\n", pokemon.Name)
	case "battle":
		fmt.Printf("You are battling %s!\n", pokemon.Name)
	case "run":
		fmt.Println("You ran away!")
	default:
		fmt.Println("Invalid choice")
	}
}
