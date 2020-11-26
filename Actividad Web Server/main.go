package main

import (
	"fmt"

	WebServer "./WebServer"
)

func main() {
	var input string
	var s WebServer.Server

	go s.Serve()

	fmt.Scanln(&input)
}
