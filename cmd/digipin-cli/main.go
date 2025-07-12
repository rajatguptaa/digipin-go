package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	digipin "github.com/rajatguptaa/digipin-go/digipin"
)

func main() {
	mode := flag.String("mode", "encode", "Mode: encode, decode, batch-encode, batch-decode")
	lat := flag.Float64("lat", 0, "Latitude (for encode)")
	lng := flag.Float64("lng", 0, "Longitude (for encode)")
	pin := flag.String("pin", "", "DIGIPIN (for decode)")
	input := flag.String("input", "", "Input CSV file for batch mode")
	output := flag.String("output", "", "Output CSV file for batch mode")
	flag.Parse()

	switch *mode {
	case "encode":
		if err := digipin.ValidateCoordinates(*lat, *lng); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		code, err := digipin.Encode(*lat, *lng)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("DIGIPIN:", code)
	case "decode":
		if err := digipin.ValidateDigiPin(*pin); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		coord, err := digipin.Decode(*pin)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Latitude: %.6f\nLongitude: %.6f\n", coord.Latitude, coord.Longitude)
	case "batch-encode":
		if *input == "" || *output == "" {
			fmt.Println("Input and output CSV files required for batch-encode.")
			os.Exit(1)
		}
		inFile, err := os.Open(*input)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			os.Exit(1)
		}
		defer inFile.Close()
		reader := csv.NewReader(inFile)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading input file:", err)
			os.Exit(1)
		}
		coords := make([]digipin.LatLng, 0, len(records))
		for _, rec := range records {
			if len(rec) < 2 {
				continue
			}
			var lat, lng float64
			fmt.Sscanf(rec[0], "%f", &lat)
			fmt.Sscanf(rec[1], "%f", &lng)
			coords = append(coords, digipin.LatLng{Latitude: lat, Longitude: lng})
		}
		pins, errs := digipin.BatchEncode(coords)
		outFile, err := os.Create(*output)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer outFile.Close()
		writer := csv.NewWriter(outFile)
		for i, pin := range pins {
			errStr := ""
			if errs[i] != nil {
				errStr = errs[i].Error()
			}
			writer.Write([]string{
				fmt.Sprintf("%f", coords[i].Latitude),
				fmt.Sprintf("%f", coords[i].Longitude),
				pin,
				errStr,
			})
		}
		writer.Flush()
		fmt.Println("Batch encode complete. Output written to", *output)
	case "batch-decode":
		if *input == "" || *output == "" {
			fmt.Println("Input and output CSV files required for batch-decode.")
			os.Exit(1)
		}
		inFile, err := os.Open(*input)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			os.Exit(1)
		}
		defer inFile.Close()
		reader := csv.NewReader(inFile)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading input file:", err)
			os.Exit(1)
		}
		pins := make([]string, 0, len(records))
		for _, rec := range records {
			if len(rec) < 1 {
				continue
			}
			pins = append(pins, strings.TrimSpace(rec[0]))
		}
		coords, errs := digipin.BatchDecode(pins)
		outFile, err := os.Create(*output)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer outFile.Close()
		writer := csv.NewWriter(outFile)
		for i, coord := range coords {
			errStr := ""
			if errs[i] != nil {
				errStr = errs[i].Error()
			}
			writer.Write([]string{
				pins[i],
				fmt.Sprintf("%f", coord.Latitude),
				fmt.Sprintf("%f", coord.Longitude),
				errStr,
			})
		}
		writer.Flush()
		fmt.Println("Batch decode complete. Output written to", *output)
	default:
		fmt.Println("Unknown mode. Use encode, decode, batch-encode, or batch-decode.")
		flag.Usage()
		os.Exit(1)
	}
}
