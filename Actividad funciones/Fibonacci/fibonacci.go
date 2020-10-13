package main

import (
	"fmt"
	"strconv"
)

func startFibonacci(iterations int) int {
	res := 0
	if iterations == 0 {
		res = 0
	} else if iterations == 1 {
		res = 1
	} else {
		res = fibo(iterations, 2, 1, 1)
	}
	return res
}

func fibo(N, it, val, prev int) int {
	if it == N {
		return val
	}
	return fibo(N, it+1, val+prev, val)
}

func main() {
	var input string

	fmt.Print("Número de iteraciones: ")
	// Read from console
	fmt.Scan(&input)

	// Scan values from input string
	N, err := strconv.Atoi(input)

	if err == nil {
		fmt.Println(startFibonacci(N))
	} else {
		fmt.Println(err)
	}
}
