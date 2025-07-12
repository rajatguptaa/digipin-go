package digipin

import (
	"fmt"
	"math"
)

// Encode converts latitude and longitude to a 10-character DIGIPIN code (format: XXX-XXX-XXXX).
// Returns an error if the coordinates are out of bounds.
func Encode(lat, lng float64) (string, error) {
	if lat < MinLat || lat > MaxLat {
		return "", fmt.Errorf("Latitude out of range")
	}
	if lng < MinLng || lng > MaxLng {
		return "", fmt.Errorf("Longitude out of range")
	}
	minLatB, maxLatB := MinLat, MaxLat
	minLngB, maxLngB := MinLng, MaxLng
	code := make([]rune, 0, 12) // 10 chars + 2 hyphens
	for level := 1; level <= 10; level++ {
		latStep := (maxLatB - minLatB) / 4.0
		lngStep := (maxLngB - minLngB) / 4.0
		row := 3 - int(math.Floor((lat-minLatB)/latStep))
		col := int(math.Floor((lng - minLngB) / lngStep))
		if row < 0 {
			row = 0
		} else if row > 3 {
			row = 3
		}
		if col < 0 {
			col = 0
		} else if col > 3 {
			col = 3
		}
		code = append(code, Grid[row][col])
		if level == 3 || level == 6 {
			code = append(code, '-')
		}
		maxLatB = minLatB + latStep*float64(4-row)
		minLatB = minLatB + latStep*float64(3-row)
		minLngB = minLngB + lngStep*float64(col)
		maxLngB = minLngB + lngStep
	}
	return string(code), nil
}
