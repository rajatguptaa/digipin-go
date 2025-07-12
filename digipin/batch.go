package digipin

import (
	"runtime"
)

// BatchEncode encodes a slice of coordinates to DIGIPINs sequentially.
// Returns a slice of codes and a slice of errors (one per input).
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

// BatchDecode decodes a slice of DIGIPINs to coordinates sequentially.
// Returns a slice of LatLng and a slice of errors (one per input).
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
// Returns a slice of codes and a slice of errors (one per input).
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
// Returns a slice of LatLng and a slice of errors (one per input).
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
