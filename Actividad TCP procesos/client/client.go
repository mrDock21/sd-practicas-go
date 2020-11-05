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

func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

func isThereError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func process(connectionDone, stopFlag chan bool, p connData.Process) {
	isStopped := false

	for !isStopped {
		select {
		case stop := <-stopFlag:
			isStopped = stop
		default: //Keep going
		}
		p.Update()
		fmt.Println(p.ToString() + "\n")
		time.Sleep(time.Millisecond * 500)
		clear()
	}
	// Send done signal to server
	conn, _ := net.Dial("tcp", ":9876")
	isThereError(gob.NewEncoder(conn).Encode(connData.ConnData{MsgType: 1, Task: p}))
	conn.Close()
	fmt.Println("[CLIENT]	Connection closed!\n")
	connectionDone <- true
}

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
	// send
	fmt.Println("[CLIENT]	Connecting...")
	if isThereError(gob.NewEncoder(c).Encode(data)) {
		c.Close()
		fmt.Println(err)
	}
	// receive
	if isThereError(gob.NewDecoder(c).Decode(&data)) {
		c.Close()
		return
	}
	c.Close()
	go process(connectionDone, stopFlag, data.Task)
}

func main() {

	var input string

	stopFlag := make(chan bool)
	connectionDoneFlag := make(chan bool)

	go connect(connectionDoneFlag, stopFlag)

	fmt.Scanln(&input)
	// Send stopping signal
	stopFlag <- true
	// Wait for connection to be clossed
	<-connectionDoneFlag
}
