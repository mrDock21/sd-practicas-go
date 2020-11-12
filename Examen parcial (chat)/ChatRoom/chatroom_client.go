package ChatRoom

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
)

const (
	DOWNLOAD_FILE_PATH = "./Download.txt"
)

type Client struct {
	Listener *net.Conn
	Username string
	Stop     *bool
}

func (c *Client) Connect(username string) bool {
	stop, isConnected := false, false
	data := Data{Cmd: Connect, Message: "Requesting connection", From: username}
	// stablish connection
	con, err := net.Dial("tcp", CONN_IP_PORT)

	if err != nil {
		fmt.Println(err)
		return isConnected
	}
	c.Listener = &con
	c.Username = username
	c.Stop = &stop
	// send initial connection petition
	gob.NewEncoder(*c.Listener).Encode(&data)
	// wait until server responds...
	gob.NewDecoder(*c.Listener).Decode(&data)

	if data.Cmd == Success {
		// Start listening incoming messages
		fmt.Println(formatToConsole(&data))
		go c.listen()
		isConnected = true
	}
	return isConnected
}

func (c *Client) listen() {
	for !*c.Stop {
		// Will wait until message is received
		data := c.receiveMessage()
		// Do something
		c.switchAction(&data)
	}
	fmt.Println("Stopping Listener...")
	(*c.Listener).Close()
}

//	Sends message to chat-room
func (c *Client) SendMessage(msg string) {
	c.sendMessage(&Data{
		Cmd:     TextStream,
		Message: msg,
		From:    c.Username,
	})
}

func (c *Client) SendFile(path string) {
	c.sendMessage(&Data{
		Cmd:     FileStream,
		Message: path,
		From:    c.Username,
	})
}

//	Sends specific data to server
func (c *Client) sendMessage(data *Data) {
	sender, err := net.Dial("tcp", CONN_IP_PORT)
	/*	Lower level function:	receives real data to connection property 	*/
	err = gob.NewEncoder(sender).Encode(data)

	if err != nil {
		fmt.Println("[ERROR WHILE SENDING]")
		fmt.Println(err)
	}
	sender.Close()
}

//	Receives data from server connection
func (c *Client) receiveMessage() Data {
	data := Data{}
	// Decode data received
	gob.NewDecoder(*c.Listener).Decode(&data)
	return data
}

func (c *Client) switchAction(data *Data) {
	switch data.Cmd {
	case Connect:
		// when another user connects (just prints)
		fmt.Println(formatToConsole(data))
	case Disconnect:
		//when another user disconnects
		fmt.Println(formatToConsole(data))
	case TextStream:
		// text received by another user
		fmt.Println(formatToConsole(data))
	case FileStream:
		// file sent by another user
		fmt.Println("---FILE RECEIVED---")
		fmt.Printf(">>> Written in path {%s}\n", DOWNLOAD_FILE_PATH)
		writeMessageToFile(data)
	default: // success is ignored
	}
}

//	Stops this client's connection and goroutines
func (c *Client) StopConnection() {
	*c.Stop = true
	c.sendMessage(&Data{
		Cmd:     Disconnect,
		Message: "Disconnect me",
		From:    c.Username,
	})
}

//	Returns a string to print in console
func formatToConsole(data *Data) string {
	return fmt.Sprintf("%s: %s", data.From, data.Message)
}

// Writes given data's message to a file
func writeMessageToFile(data *Data) {
	// Write to file
	err := ioutil.WriteFile(DOWNLOAD_FILE_PATH, []byte(data.Message), 0644)

	if err != nil {
		fmt.Println(err)
	}
}
