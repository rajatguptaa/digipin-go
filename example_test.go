package digipin_test

import (
	"fmt"

	digipin "github.com/rajatgupta/digipin-go/digipin"
)

func ExampleEncode() {
	pin, err := digipin.Encode(28.6139, 77.2090)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(pin)
	// Output: 39J-438-TJC7
}

func ExampleDecode() {
	coord, err := digipin.Decode("39J-438-TJC7")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%.6f, %.6f\n", coord.Latitude, coord.Longitude)
	// Output: 28.613901, 77.208998
}

func ExampleGetDistance() {
	d1, err := digipin.GetDistance("39J-438-TJC7", "4FK-595-8823")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%.0f\n", d1)
	// Output: 1148096
}

func ExampleOrderByDistance() {
	pins := []string{"4FK-595-8823", "4PJ-766-C924"}
	ordered, err := digipin.OrderByDistance("39J-438-TJC7", pins)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(ordered)
	// Output: [4FK-595-8823 4PJ-766-C924]
}

func ExampleFindNearest() {
	pins := []string{"4FK-595-8823", "4PJ-766-C924", "422-5C2-LTTF"}
	nearest, err := digipin.FindNearest("39J-438-TJC7", pins)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(nearest)
	// Output: 4FK-595-8823
}
