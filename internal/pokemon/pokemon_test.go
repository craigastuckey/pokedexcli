package pokemon

import (
	"testing"
)

func TestGetPokemon(t *testing.T) {
	pokemonName := "pikachu"
	_, err := GetPokemon(pokemonName)
	if err != nil { // Expecting no error for a valid pokemon
		t.Fatalf("Expected no error, got %v", err)
	}

	pokemonName = "invalidpokemon"
	_, err = GetPokemon(pokemonName)
	if err == nil { // Expecting an error for an invalid pokemon
		t.Fatalf("Expected error for invalid pokemon, got nil")
	}
}
