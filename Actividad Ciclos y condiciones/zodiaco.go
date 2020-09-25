package main

import (
	"fmt"
	"strconv"
)

func main() {
	var dayInput, monthInput, res string

	// Read from console
	fmt.Scan(&dayInput)
	fmt.Scan(&monthInput)

	// Scan values from input string
	day, _ := strconv.Atoi(dayInput)
	month, _ := strconv.Atoi(monthInput)

	switch {
	case month == 12 && day >= 22 || month == 1 && day <= 20:
		res = "capricornio"
	case month == 11 && day >= 23 || month == 12 && day <= 21:
		res = "sagitario"
	case month == 1 && day >= 21 || month == 2 && day <= 18:
		res = "acuario"
	case month == 10 && day >= 23 || month == 11 && day >= 22:
		res = "escorpio"
	case month == 9 && day >= 23 || month == 10 && day <= 22:
		res = "libra"
	case month == 8 && day >= 23 || month == 9 && day <= 22:
		res = "virgo"
	case month == 7 && day >= 23 || month == 8 && day <= 22:
		res = "leo"
	case month == 6 && day >= 22 || month == 7 && day <= 22:
		res = "cancer"
	case month == 5 && day >= 21 || month == 6 && day <= 21:
		res = "geminis"
	case month == 4 && day >= 21 || month == 5 && day <= 20:
		res = "tauro"
	case month == 3 && day >= 21 || month == 4 && day <= 20:
		res = "aries"
	case month == 2 && day >= 19 || month == 3 && day <= 20:
		res = "piscis"
	}
	fmt.Println(res)
}
