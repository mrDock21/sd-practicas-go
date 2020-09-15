package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input string
	var fahrenheit int = 1

	fmt.Print(`
	+---------FAHRENHEIT A CELCIUS----------+
	  Ingrese los grados fahrenheit: `)

	// Read from console
	_, err := fmt.Scan(&input)

	if err == nil {
		// Scan values from input string
		integer, atoi_err := strconv.Atoi(input)
		if atoi_err == nil {
			fahrenheit = integer
			fmt.Printf("\t\tGrados Celcius = 5/9 x (%d - 32) = %.3f\n", fahrenheit, (5.0/9.0)*float32(fahrenheit-32))
		} else {
			fmt.Println("\t\tValor inválido!")
		}
	} else {
		fmt.Println(err)
	}
}
