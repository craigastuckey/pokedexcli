package location

import (
	"testing"
)

func TestGetLocationArea(t *testing.T) {
	locationArea, err := GetLocationArea("https://pokeapi.co/api/v2/location-area/1/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if locationArea.Name != "canalave-city-area" {
		t.Fatalf("Expected name to be canalave-city-area, got %s", locationArea.Name)
	}

	locationArea, err = GetLocationArea("https://pokeapi.co/api/v2/location-area/0/")
	if err == nil {
		t.Fatalf("Expected error for invalid location area, got nil")
	}
}
