package main

import (
	"fmt"
	"strconv"
)

func factorial(number int) float32 {
	var res float32 = 0

	if number > 1 {
		aux := float32(number)
		for i := number - 1; i > 0; i-- {
			res += aux * float32(i)
			aux = float32(i)
		}
	} else {
		res = 1
	}

	return res
}

func main() {
	var input string
	var res float32 = 0

	fmt.Print("Número de iteraciones: ")
	// Read from console
	fmt.Scan(&input)

	// Scan values from input string
	N, _ := strconv.Atoi(input)

	for i := 0; i < N; i++ {
		res += float32(1) / factorial(i)
	}

	fmt.Println(res)
}
