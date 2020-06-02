package main2

import (
	"fmt"
	"math/rand"

	"golang.org/x/tour/pic"
)

// Pic ddd
func Pic(dx, dy int) [][]uint8 {

	data := make([][]uint8, dy)
	// fmt.Printf("%T %v\n", data, data)

	for i := range data {
		data[i] = make([]uint8, dx)
		// fmt.Printf("index %v %T %v\n", i, raw, raw)

		for j := range data[i] {
			data[i][j] = uint8(rand.Intn(255))
		}
	}
	fmt.Printf("%T %v\n", data, data)
	return data
}

func main() {
	pic.Show(Pic)
}
