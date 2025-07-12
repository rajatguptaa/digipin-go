package digipin

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"encoding/csv"
	"os"

	"runtime"

	lru "github.com/hashicorp/golang-lru/v2"
)

// DIGIPIN character set (16 characters)
var charset = []rune{'F', 'C', '9', '8', 'J', '3', '2', '7', 'K', '4', '5', '6', 'L', 'M', 'P', 'T'}

// DIGIPIN 4x4 grid
var digipinGrid = [4][4]rune{
	{'F', 'C', '9', '8'},
	{'J', '3', '2', '7'},
	{'K', '4', '5', '6'},
	{'L', 'M', 'P', 'T'},
}

// Bounds for India
const (
	minLat = 2.5
	maxLat = 38.5
	minLng = 63.5
	maxLng = 99.5
)

// Result struct for decoded coordinates
type LatLng struct {
	Latitude  float64
	Longitude float64
}

// LRU cache for coordinate->DIGIPIN and DIGIPIN->coordinate lookups
var (
	cacheOnce   sync.Once
	encodeCache *lru.Cache[string, string]
	decodeCache *lru.Cache[string, LatLng]
)

func initCaches() {
	cacheOnce.Do(func() {
		encodeCache, _ = lru.New[string, string](10000)
		decodeCache, _ = lru.New[string, LatLng](10000)
	})
}

// GetCachedEncode returns cached DIGIPIN for lat/lng if present
func GetCachedEncode(lat, lng float64) (string, bool) {
	initCaches()
	key := fmt.Sprintf("%f,%f", lat, lng)
	val, ok := encodeCache.Get(key)
	return val, ok
}

// SetCachedEncode stores DIGIPIN for lat/lng
func SetCachedEncode(lat, lng float64, pin string) {
	initCaches()
	key := fmt.Sprintf("%f,%f", lat, lng)
	encodeCache.Add(key, pin)
}

// GetCachedDecode returns cached LatLng for DIGIPIN if present
func GetCachedDecode(pin string) (LatLng, bool) {
	initCaches()
	val, ok := decodeCache.Get(pin)
	return val, ok
}

// SetCachedDecode stores LatLng for DIGIPIN
func SetCachedDecode(pin string, coord LatLng) {
	initCaches()
	decodeCache.Add(pin, coord)
}

// Clear all caches
func ClearCache() {
	initCaches()
	encodeCache.Purge()
	decodeCache.Purge()
}

// Encode coordinates to DIGIPIN (XXX-XXX-XXXX)
func Encode(lat, lng float64) (string, error) {
	if lat < minLat || lat > maxLat {
		return "", fmt.Errorf("Latitude out of range")
	}
	if lng < minLng || lng > maxLng {
		return "", fmt.Errorf("Longitude out of range")
	}
	minLatB, maxLatB := minLat, maxLat
	minLngB, maxLngB := minLng, maxLng
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
		code = append(code, digipinGrid[row][col])
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

// Decode DIGIPIN to coordinates (center of cell)
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
	minLatB, maxLatB := minLat, maxLat
	minLngB, maxLngB := minLng, maxLng
	for _, char := range clean {
		found := false
		var row, col int
		for r := 0; r < 4; r++ {
			for c := 0; c < 4; c++ {
				if digipinGrid[r][c] == char {
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

// Validate latitude and longitude bounds
func ValidateCoordinates(lat, lng float64) error {
	if math.IsNaN(lat) || math.IsNaN(lng) {
		return fmt.Errorf("Invalid coordinates: must be numbers")
	}
	if lat < minLat || lat > maxLat {
		return fmt.Errorf("Latitude out of range")
	}
	if lng < minLng || lng > maxLng {
		return fmt.Errorf("Longitude out of range")
	}
	return nil
}

// Validate DIGIPIN format (10 chars, valid characters)
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
				if digipinGrid[r][c] == char {
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

// BatchEncode encodes a slice of coordinates to DIGIPINs
func BatchEncode(coords []LatLng) ([]string, []error) {
	results := make([]string, len(coords))
	errs := make([]error, len(coords))
	for i, c := range coords {
		if err := ValidateCoordinates(c.Latitude, c.Longitude); err != nil {
			results[i] = ""
			errs[i] = err
			continue
		}
		pin, err := Encode(c.Latitude, c.Longitude)
		results[i] = pin
		errs[i] = err
	}
	return results, errs
}

// BatchDecode decodes a slice of DIGIPINs to coordinates
func BatchDecode(pins []string) ([]LatLng, []error) {
	results := make([]LatLng, len(pins))
	errs := make([]error, len(pins))
	for i, pin := range pins {
		if err := ValidateDigiPin(pin); err != nil {
			results[i] = LatLng{}
			errs[i] = err
			continue
		}
		coord, err := Decode(pin)
		results[i] = coord
		errs[i] = err
	}
	return results, errs
}

// BatchEncodeConcurrent encodes a slice of coordinates to DIGIPINs concurrently using a worker pool.
func BatchEncodeConcurrent(coords []LatLng) ([]string, []error) {
	type result struct {
		idx int
		pin string
		err error
	}
	results := make([]string, len(coords))
	errs := make([]error, len(coords))
	ch := make(chan result, len(coords))
	workers := runtime.NumCPU()
	sem := make(chan struct{}, workers)
	for i, c := range coords {
		sem <- struct{}{} // acquire
		go func(i int, c LatLng) {
			defer func() { <-sem }() // release
			if err := ValidateCoordinates(c.Latitude, c.Longitude); err != nil {
				ch <- result{i, "", err}
				return
			}
			pin, err := Encode(c.Latitude, c.Longitude)
			ch <- result{i, pin, err}
		}(i, c)
	}
	for i := 0; i < len(coords); i++ {
		r := <-ch
		results[r.idx] = r.pin
		errs[r.idx] = r.err
	}
	return results, errs
}

// BatchDecodeConcurrent decodes a slice of DIGIPINs to coordinates concurrently using a worker pool.
func BatchDecodeConcurrent(pins []string) ([]LatLng, []error) {
	type result struct {
		idx   int
		coord LatLng
		err   error
	}
	results := make([]LatLng, len(pins))
	errs := make([]error, len(pins))
	ch := make(chan result, len(pins))
	workers := runtime.NumCPU()
	sem := make(chan struct{}, workers)
	for i, pin := range pins {
		sem <- struct{}{} // acquire
		go func(i int, pin string) {
			defer func() { <-sem }() // release
			if err := ValidateDigiPin(pin); err != nil {
				ch <- result{i, LatLng{}, err}
				return
			}
			coord, err := Decode(pin)
			ch <- result{i, coord, err}
		}(i, pin)
	}
	for i := 0; i < len(pins); i++ {
		r := <-ch
		results[r.idx] = r.coord
		errs[r.idx] = r.err
	}
	return results, errs
}

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
