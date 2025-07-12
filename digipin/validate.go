package digipin

import (
	"fmt"
	"math"
)

// ValidateCoordinates checks if latitude and longitude are within valid bounds and are numbers.
// Returns an error if invalid.
func ValidateCoordinates(lat, lng float64) error {
	if math.IsNaN(lat) || math.IsNaN(lng) {
		return fmt.Errorf("Invalid coordinates: must be numbers")
	}
	if lat < MinLat || lat > MaxLat {
		return fmt.Errorf("Latitude out of range")
	}
	if lng < MinLng || lng > MaxLng {
		return fmt.Errorf("Longitude out of range")
	}
	return nil
}

// ValidateDigiPin checks if the DIGIPIN is 10 valid characters (excluding hyphens) and only uses allowed characters.
// Returns an error if invalid.
func ValidateDigiPin(pin string) error {
	clean := make([]rune, 0, 10)
	for _, r := range pin {
		if r == '-' {
			continue
		}
		clean = append(clean, r)
	}
	if len(clean) != 10 {
		return fmt.Errorf("Invalid DIGIPIN: must be 10 characters (excluding hyphens)")
	}
	for _, char := range clean {
		found := false
		for r := 0; r < 4; r++ {
			for c := 0; c < 4; c++ {
				if Grid[r][c] == char {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return fmt.Errorf("Invalid character in DIGIPIN: %c", char)
		}
	}
	return nil
}
