package main

import (
	"fmt"
)

func findBiggest(numbers ...int) int {
	biggest := -1

	for _, v := range numbers {
		if v > biggest {
			biggest = v
		}
	}
	return biggest
}

func main() {
	var N, number int = 0, 0
	var slice []int

	fmt.Print("Número de elementos: ")
	// Read from console
	fmt.Scanln(&N)
	// Read each element
	for i := 0; i < N; i++ {
		fmt.Scanln(&number)
		slice = append(slice, number)
	}
	fmt.Printf("Número más grande: %d\n", findBiggest(slice...))
}
