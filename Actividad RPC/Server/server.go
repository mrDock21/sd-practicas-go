package main

import (
	"fmt"

	ServerRPC "../ServerRPC"
)

func main() {
	var input string
	server := ServerRPC.Server{}

	go server.Serve()

	fmt.Scanln(&input)
}
