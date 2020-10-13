package main

import (
	"fmt"
)

func oddNumberBuilder() func() int {
	i := 1

	return func() int {
		var odd = i
		i += 2
		return odd
	}
}

func main() {
	var nextOdd = oddNumberBuilder()

	fmt.Println(nextOdd())
	fmt.Println(nextOdd())
	fmt.Println(nextOdd())
	fmt.Println(nextOdd())
	fmt.Println(nextOdd())
}
