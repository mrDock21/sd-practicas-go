package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	ChatRoom "../ChatRoom"
)

//	Clears console
func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

//	Main client's loop to interact with the chat
func clientActions(client *ChatRoom.Client) {
	var input string

	isChatting := true

	for isChatting {

		clear()
		fmt.Println("[NOTE: Downloaded files are saved in current path]")
		fmt.Println("+-------OPTIONS-------+")
		fmt.Println(" a) Send message")
		fmt.Println(" b) Send file")
		fmt.Println(" c) Exit")

		fmt.Printf("\n+---CHAT---+\n%s\n", client.ChatHistory)

		fmt.Println(">>>")
		fmt.Scanln(&input)

		switch {
		case input == "a":
			// Send text message
			client.SendMessage(readMessage())
		case input == "b":
			fileName := readMessage()
			// Send file's text
			contents, success := tryReadFile(fileName)
			if success {
				client.SendFile(contents, fileName)
			}
		}
		isChatting = input != "c"
	}
}

//	Reads user input from console
func readMessage() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

//	Will try to read file's contents
//	Returns false if failure
func tryReadFile(path string) (string, bool) {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return "", false
	}
	return string(contents), true
}

func main() {
	var input string
	client := ChatRoom.Client{}

	fmt.Println("----CHAT----")
	fmt.Print("Type your username: ")
	fmt.Scanln(&input)

	if client.Connect(input) {
		// main loop of actions...
		clientActions(&client)
		// disconnect when done
		client.StopConnection()
	}
}
