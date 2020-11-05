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

//	Clears console window
func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

// Updates all processes' process time
func updateTimers(processes *[]connData.Process) {
	for i := 0; i < len(*processes); i++ {
		(*processes)[i].Update()
		fmt.Println((*processes)[i].ToString())
	}
	fmt.Println("")
}

//	Main process loop goroutine. Stops until stopFlag receives "true" signal
func process(stopFlag chan bool, processes *[]connData.Process) {
	isProcessInterrupted := false
	for {
		// Listen channels
		select {
		case stop := <-stopFlag:
			isProcessInterrupted = stop
		default: // Keep running
		}

		if !isProcessInterrupted {
			// Only update processes times while in safe mode
			updateTimers(processes)
			time.Sleep(time.Millisecond * 500)
			clear()
		}
	}
}

//	Sends process to client at index 0 of processes'
//	slice and removes it from main slice
func sendProcessToClient(connection net.Conn, processes *[]connData.Process) {
	// Make copy of first process
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
	// truncate slice (removing last element)
	*processes = (*processes)[:len(*processes)-1]
	// send process to client
	gob.NewEncoder(connection).Encode(response)
}

//	Main server loop. Listens for client's messages
func serve() {
	// Initial hardcoded processes
	processes := []connData.Process{
		connData.Process{Id: 1, Time: 0},
		connData.Process{Id: 2, Time: 0},
		connData.Process{Id: 3, Time: 0},
		connData.Process{Id: 4, Time: 0},
		connData.Process{Id: 5, Time: 0},
	}

	stopFlag := make(chan bool)
	// Start server
	server, err := net.Listen("tcp", ":9876")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("[SERVER]	Listening...")
	// Start counting processes' times
	go process(stopFlag, &processes)

	for {
		// wait for client's connection
		connection, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(connection, stopFlag, &processes)
	}
}

//	Decides which operation to make according to client's message
func handleClient(connection net.Conn, interruptProcess chan bool, processes *[]connData.Process) {
	// GUARD (avoids deleting processes while been processed)
	interruptProcess <- true

	var data connData.ConnData

	gob.NewDecoder(connection).Decode(&data)

	if data.MsgType > 0 {
		// Receiving process
		fmt.Println("[SERVER]	RECEIVES PROCESS")
		*processes = append(*processes, data.Task)
	} else {
		// Send process
		fmt.Println("[SERVER]	SENDS PROCESS")
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
