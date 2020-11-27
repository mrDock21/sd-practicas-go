package main

import (
	"fmt"

	WebServer "./WebServer"
)

func main() {
	var input string
	var s WebServer.Server

	go s.Serve()
	fmt.Println("[SERVER]	Running")
	fmt.Scanln(&input)
}
