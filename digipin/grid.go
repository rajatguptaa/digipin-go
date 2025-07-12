package digipin

import (
	"encoding/csv"
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
