package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input string
	var side int = 1

	fmt.Print(`
	+---------AREA DEL CUADRADO----------+
	  Ingrese el lado del cuadrado: `)

	// Read from console
	_, err := fmt.Scan(&input)

	if err == nil {
		// Scan values from input string
		integer, atoi_err := strconv.Atoi(input)
		if atoi_err == nil {
			side = integer
			fmt.Printf("\t\tEl Área = %d x %d = %d\n", side, side, side*side)
		} else {
			fmt.Println("\t\tValor inválido!")
		}
	} else {
		fmt.Println(err)
	}
}
