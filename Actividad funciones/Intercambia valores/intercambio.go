package main

import (
	"fmt"
)

func swap(a, b *int) {
	var aux = *a

	*a = *b
	*b = aux
}

func main() {
	var a, b int = 0, 0

	fmt.Print("Ingrese los valores (a, b):")
	fmt.Scanf("%d, %d ", &a, &b)

	swap(&a, &b)

	fmt.Printf("Valores intercambiados (a, b): %d, %d\n", a, b)
}
