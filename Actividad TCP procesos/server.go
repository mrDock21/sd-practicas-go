package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	connData "./ConnData"
)

func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

func updateTimers(processes *[]connData.Process) {
	for i := 0; i < len(*processes); i++ {
		(*processes)[i].Update()
		fmt.Println((*processes)[i].ToString())
	}
	fmt.Println("")
}

func process(stopFlag chan bool, processes *[]connData.Process) {
	isProcessInterrupted := false
	for {

		select {
		case stop := <-stopFlag:
			isProcessInterrupted = stop
		default: // Keep running
		}

		if !isProcessInterrupted {
			updateTimers(processes)
			time.Sleep(time.Millisecond * 500)
			clear()
		}
	}
}

func sendProcessToClient(connection net.Conn, processes *[]connData.Process) {
	sentProcess := connData.Process{
		Id: (*processes)[0].Id, Time: (*processes)[0].Time,
	}
	response := connData.ConnData{
		MsgType: 0,
		Task:    sentProcess,
	}
	// Move all elements 1 to the left (deleting first)
	for i := 0; i < len(*processes)-1; i++ {
		(*processes)[i] = (*processes)[i+1]
	}
	// truncate
	*processes = (*processes)[:len(*processes)-1]
	// send process to client
	gob.NewEncoder(connection).Encode(response)
}

func serve() {
	processes := []connData.Process{
		connData.Process{Id: 1, Time: 0},
		connData.Process{Id: 2, Time: 0},
		connData.Process{Id: 3, Time: 0},
		connData.Process{Id: 4, Time: 0},
		connData.Process{Id: 5, Time: 0},
	}

	stopFlag := make(chan bool)

	server, err := net.Listen("tcp", ":9876")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("[SERVER]	Listening...")

	go process(stopFlag, &processes)

	for {
		connection, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(connection, stopFlag, &processes)
	}
}

func handleClient(connection net.Conn, interruptProcess chan bool, processes *[]connData.Process) {
	// GUARD
	interruptProcess <- true

	var data connData.ConnData

	gob.NewDecoder(connection).Decode(&data)

	if data.MsgType > 0 {
		fmt.Println("[SERVER]	RECEIVES PROCESS")
		// Receiving process
		*processes = append(*processes, data.Task)
	} else {
		fmt.Println("[SERVER]	SENDS PROCESS")
		// Send process
		sendProcessToClient(connection, processes)
	}

	interruptProcess <- false
}

func main() {
	var input string

	go serve()

	fmt.Print("[SERVER] Press enter to stop server...")
	fmt.Scanln(&input)
}
