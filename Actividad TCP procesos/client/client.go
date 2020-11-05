package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	connData "../ConnData"
)

//	Clears console
func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

//	General printer for error messages
func isThereError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

//	Process individual process given by server. Stops after receiving stopping signal (channel)
func process(connectionDone, stopFlag chan bool, p connData.Process) {
	isStopped := false

	for !isStopped {
		// Listen for channels
		select {
		case stop := <-stopFlag:
			isStopped = stop
		default: //Keep going
		}
		// Process updating
		p.Update()
		fmt.Println(p.ToString() + "\n")
		time.Sleep(time.Millisecond * 500)
		clear()
	}
	// Send done signal to server
	conn, _ := net.Dial("tcp", ":9876")
	// Give back process to server
	isThereError(gob.NewEncoder(conn).Encode(connData.ConnData{MsgType: 1, Task: p}))
	conn.Close()
	fmt.Println("[CLIENT]	Connection closed!\n")
	// Allow main to finish
	connectionDone <- true
}

//	Connects to server and receives process from it.
func connect(connectionDone, stopFlag chan bool) {
	c, err := net.Dial("tcp", ":9876")

	if isThereError(err) {
		return
	}

	data := connData.ConnData{
		MsgType: -1,
		Task: connData.Process{
			Id: 0, Time: 0,
		},
	}
	// send initial message to server
	fmt.Println("[CLIENT]	Connecting...")
	if isThereError(gob.NewEncoder(c).Encode(data)) {
		c.Close()
		fmt.Println(err)
	}
	// receive process from server
	if isThereError(gob.NewDecoder(c).Decode(&data)) {
		c.Close()
		return
	}
	c.Close()
	// start computing given process
	go process(connectionDone, stopFlag, data.Task)
}

func main() {

	var input string

	stopFlag := make(chan bool)
	connectionDoneFlag := make(chan bool)

	go connect(connectionDoneFlag, stopFlag)
	// wait until user presses "enter"
	fmt.Scanln(&input)
	// Send stopping signal
	stopFlag <- true
	// Wait for connection to be clossed
	<-connectionDoneFlag
}
