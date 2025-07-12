# digipin-go

A comprehensive Go library and CLI for encoding and decoding Indian geographic coordinates into DIGIPIN format (Indian Postal Digital PIN system).

## 🚀 Features
- **Coordinate to DIGIPIN:** Convert lat/lng to 10-character code
- **DIGIPIN to Coordinates:** Reverse lookup with high precision
- **Validation:** Built-in coordinate bounds and format checking
- **Batch Processing:** Encode/decode multiple coordinates or pins at once (sequential and concurrent)
- **LRU Caching:** Fast repeated lookups
- **Geospatial Utilities:** Distance, nearest, order by distance
- **Grid Generation:** Create offline coordinate grids (CSV)
- **CLI Tool:** Command-line interface for all major features
- **Modular Go Package:** Clean, idiomatic, and production-ready

## 📦 Installation

```
go get github.com/rajatgupta/digipin-go/digipin
```

## 🛠️ CLI Usage

Build and run the CLI:

```
go run ./cmd/digipin-cli -mode encode -lat 28.6139 -lng 77.2090
# Output: DIGIPIN: 39J-438-TJC7

go run ./cmd/digipin-cli -mode decode -pin 39J-438-TJC7
# Output: Latitude: 28.613901
#         Longitude: 77.208998
```

Batch encode/decode:

```
go run ./cmd/digipin-cli -mode batch-encode -input example_batch_encode.csv -output output_batch_encode.csv
go run ./cmd/digipin-cli -mode batch-decode -input input_batch_decode.csv -output output_batch_decode.csv
```

## 📚 Library Usage

```go
import digipin "github.com/rajatgupta/digipin-go/digipin"

// Encode coordinates
pin, err := digipin.Encode(28.6139, 77.2090) // "39J-438-TJC7"

// Decode DIGIPIN
coord, err := digipin.Decode("39J-438-TJC7") // {28.613901, 77.208998}

// Batch encode
coords := []digipin.LatLng{
    {Latitude: 28.6139, Longitude: 77.2090},
    {Latitude: 19.0760, Longitude: 72.8777},
}
pins, errs := digipin.BatchEncode(coords)

// Geospatial utilities
km, _ := digipin.GetDistance("39J-438-TJC7", "4FK-595-8823")
nearest, _ := digipin.FindNearest("39J-438-TJC7", []string{"4FK-595-8823", "4PJ-766-C924"})

// Grid generation
err := digipin.GenerateGrid(20, 70, 21, 71, 0.1, "grid.csv")
```

## 🧪 Examples & Tests

Run all examples:
```
go test -v example_test.go
```

## 🤝 Contributing

Contributions are welcome! Please open issues or pull requests for improvements, bug fixes, or new features.

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.

---

Made with ❤️ for the Indian developer community. 