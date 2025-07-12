package digipin

import (
	"fmt"
)

// Decode converts a 10-character DIGIPIN code (format: XXX-XXX-XXXX) back to the center latitude and longitude of the corresponding cell.
// Returns an error if the code is invalid or out of bounds.
func Decode(pin string) (LatLng, error) {
	clean := make([]rune, 0, 10)
	for _, r := range pin {
		if r == '-' {
			continue
		}
		clean = append(clean, r)
	}
	if len(clean) != 10 {
		return LatLng{}, fmt.Errorf("Invalid DIGIPIN: must be 10 characters (excluding hyphens)")
	}
	minLatB, maxLatB := MinLat, MaxLat
	minLngB, maxLngB := MinLng, MaxLng
	for _, char := range clean {
		found := false
		var row, col int
		for r := 0; r < 4; r++ {
			for c := 0; c < 4; c++ {
				if Grid[r][c] == char {
					row, col = r, c
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return LatLng{}, fmt.Errorf("Invalid character in DIGIPIN: %c", char)
		}
		latStep := (maxLatB - minLatB) / 4.0
		lngStep := (maxLngB - minLngB) / 4.0
		newMaxLat := minLatB + latStep*float64(4-row)
		newMinLat := minLatB + latStep*float64(3-row)
		newMinLng := minLngB + lngStep*float64(col)
		newMaxLng := newMinLng + lngStep
		minLatB, maxLatB, minLngB, maxLngB = newMinLat, newMaxLat, newMinLng, newMaxLng
	}
	return LatLng{
		Latitude:  (minLatB + maxLatB) / 2.0,
		Longitude: (minLngB + maxLngB) / 2.0,
	}, nil
}
