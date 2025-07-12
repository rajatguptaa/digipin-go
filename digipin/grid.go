package digipin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// GenerateGrid generates a grid of DIGIPINs for the given bounding box and step size, and writes to a CSV file.
func GenerateGrid(minLat, minLng, maxLat, maxLng, step float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write header
	writer.Write([]string{"latitude", "longitude", "digipin", "error"})
	for lat := minLat; lat <= maxLat; lat += step {
		for lng := minLng; lng <= maxLng; lng += step {
			pin, err := Encode(lat, lng)
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			writer.Write([]string{
				fmt.Sprintf("%.6f", lat),
				fmt.Sprintf("%.6f", lng),
				pin,
				errStr,
			})
		}
	}
	return nil
}

// GridEntry represents a single entry in the grid for JSON output.
type GridEntry struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DigiPin   string  `json:"digipin"`
	Error     string  `json:"error"`
}

// GenerateGridJSON generates a grid of DIGIPINs and writes to a JSON file as an array of objects.
func GenerateGridJSON(minLat, minLng, maxLat, maxLng, step float64, filename string) error {
	entries := make([]GridEntry, 0)
	for lat := minLat; lat <= maxLat; lat += step {
		for lng := minLng; lng <= maxLng; lng += step {
			pin, err := Encode(lat, lng)
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			entries = append(entries, GridEntry{
				Latitude:  lat,
				Longitude: lng,
				DigiPin:   pin,
				Error:     errStr,
			})
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}
