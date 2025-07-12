package digipin

// LatLng represents a latitude/longitude pair.
type LatLng struct {
	Latitude  float64 // Latitude in decimal degrees
	Longitude float64 // Longitude in decimal degrees
}

// Charset is the DIGIPIN character set (16 characters).
var Charset = []rune{'F', 'C', '9', '8', 'J', '3', '2', '7', 'K', '4', '5', '6', 'L', 'M', 'P', 'T'}

// Grid is the DIGIPIN 4x4 character grid.
var Grid = [4][4]rune{
	{'F', 'C', '9', '8'},
	{'J', '3', '2', '7'},
	{'K', '4', '5', '6'},
	{'L', 'M', 'P', 'T'},
}

// Bounds for valid Indian coordinates.
const (
	MinLat = 2.5  // Minimum latitude (degrees)
	MaxLat = 38.5 // Maximum latitude (degrees)
	MinLng = 63.5 // Minimum longitude (degrees)
	MaxLng = 99.5 // Maximum longitude (degrees)
)
