package digipin

import (
	"fmt"
	"math"
	"sort"
)

// HaversineDistance calculates the distance in meters between two coordinates
func HaversineDistance(a, b LatLng) float64 {
	const R = 6371000.0 // Earth radius in meters
	lat1 := a.Latitude * math.Pi / 180.0
	lat2 := b.Latitude * math.Pi / 180.0
	dlat := (b.Latitude - a.Latitude) * math.Pi / 180.0
	dlng := (b.Longitude - a.Longitude) * math.Pi / 180.0
	sinDlat := math.Sin(dlat / 2)
	sinDlng := math.Sin(dlng / 2)
	aVal := sinDlat*sinDlat + math.Cos(lat1)*math.Cos(lat2)*sinDlng*sinDlng
	c := 2 * math.Atan2(math.Sqrt(aVal), math.Sqrt(1-aVal))
	return R * c
}

// GetDistance returns the Haversine distance between two DIGIPINs (meters)
func GetDistance(pinA, pinB string) (float64, error) {
	if err := ValidateDigiPin(pinA); err != nil {
		return 0, err
	}
	if err := ValidateDigiPin(pinB); err != nil {
		return 0, err
	}
	a, err := Decode(pinA)
	if err != nil {
		return 0, err
	}
	b, err := Decode(pinB)
	if err != nil {
		return 0, err
	}
	return HaversineDistance(a, b), nil
}

// GetPreciseDistance is an alias for GetDistance (Vincenty can be added later)
func GetPreciseDistance(pinA, pinB string) (float64, error) {
	return GetDistance(pinA, pinB)
}

// OrderByDistance sorts pins by distance from a reference pin
func OrderByDistance(referencePin string, pins []string) ([]string, error) {
	if err := ValidateDigiPin(referencePin); err != nil {
		return nil, err
	}
	ref, err := Decode(referencePin)
	if err != nil {
		return nil, err
	}
	type pinDist struct {
		pin  string
		dist float64
	}
	list := make([]pinDist, 0, len(pins))
	for _, pin := range pins {
		if err := ValidateDigiPin(pin); err != nil {
			return nil, err
		}
		coord, err := Decode(pin)
		if err != nil {
			return nil, err
		}
		dist := HaversineDistance(ref, coord)
		list = append(list, pinDist{pin, dist})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].dist < list[j].dist
	})
	result := make([]string, len(list))
	for i, pd := range list {
		result[i] = pd.pin
	}
	return result, nil
}

// FindNearest returns the nearest pin to the reference pin
func FindNearest(referencePin string, pins []string) (string, error) {
	ordered, err := OrderByDistance(referencePin, pins)
	if err != nil {
		return "", err
	}
	if len(ordered) == 0 {
		return "", fmt.Errorf("No pins provided")
	}
	return ordered[0], nil
}
