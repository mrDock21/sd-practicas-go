package main

import (
	"fmt"

	ChatRoom "../ChatRoom"
)

//	Main client's loop to interact with the chat
func serverActions(server *ChatRoom.Server) {
	var input string

	isServerRunning := true

	for isServerRunning {
		fmt.Println("\n+-------OPTIONS-------+")
		fmt.Println(" a) Show Chat history")
		fmt.Println(" b) Save chat history")
		fmt.Println(" c) Stop server")

		fmt.Println("Option:")
		fmt.Scanln(&input)

		switch {
		case input == "a":
			// Send text message
			fmt.Println("----CHAT-HISTORY----")
			fmt.Printf("\n%s\n", server.ChatHistory)
		case input == "b":
			fmt.Println("Saved to ChatHistory.txt")
			server.BackupChat()
		}
		isServerRunning = input != "c"
	}
}

func main() {

	server := ChatRoom.Server{}

	// start chat-server
	server.Start()
	// show menu
	serverActions(&server)
}
