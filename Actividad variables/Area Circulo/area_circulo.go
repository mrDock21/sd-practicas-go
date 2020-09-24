package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input string
	var radius int = 1
	const PI = 3.141592

	fmt.Print(`
	+---------AREA DEL CÍRCULO----------+
	  Ingrese la radio del círculo: `)

	// Read from console
	_, err := fmt.Scan(&input)

	if err == nil {
		// Scan values from input string
		integer, atoi_err := strconv.Atoi(input)
		if atoi_err == nil {
			radius = integer
			fmt.Printf("\t\tÁrea = PI x %d x %d = %.3f\n", radius, radius, PI*float32(radius)*float32(radius))
		} else {
			fmt.Println("\t\tValor inválido!")
		}
	} else {
		fmt.Println(err)
	}
}
