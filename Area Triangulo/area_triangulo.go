package main

import (
	"fmt"
)

func psedoTernary(expression bool, trueRes string, falseRes string) string {
	res := ""
	if expression {
		res = trueRes
	} else {
		res = falseRes
	}
	return res
}

func main() {
	var base, height int = 1, 1

	fmt.Print(`
	+---------AREA DEL TRIÁNGULO----------+
	  Ingrese la base y altura (separados por espacio): `)

	// Read from console
	num_args, err := fmt.Scanf("%d %d", &base, &height)

	if err == nil && num_args == 2 {
		fmt.Printf("\t\tÁrea = (%d x %d) / 2 = %.3f\n", base, height, float32(base)*float32(height)/2.0)
	} else {
		fmt.Println(err)
		fmt.Println("\n" + psedoTernary(num_args < 2, "Debe ingresar la base y la altura", ""))
	}
}
