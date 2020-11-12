package main

import (
	"fmt"

	ChatRoom "../ChatRoom"
)

func main() {
	var input string
	server := ChatRoom.Server{}

	// start chat-server
	server.Start()
	// wait until user presses enter to stop server
	fmt.Scanln(&input)
}
