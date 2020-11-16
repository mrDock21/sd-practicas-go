package ChatRoom

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

const (
	DOWNLOAD_FILE_PATH = "./"
)

type Client struct {
	Listener *net.Conn
	Username string
	Stop     *bool
	// to shwo all messages
	ChatHistory string
	// used when someone sends a file
	fileName string
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
	c.ChatHistory += fmt.Sprintf("You: {%s}\n", msg)
}

func (c *Client) SendFile(fileContents, fileName string) {
	// send notification
	c.sendMessage(&Data{
		Cmd:     PrepareFile,
		Message: fileName,
		From:    c.Username,
	})
	// force the sleep to make sure name comes first
	time.Sleep(time.Millisecond * 1000)
	// send stream
	c.sendMessage(&Data{
		Cmd:     FileStream,
		Message: fileContents,
		From:    c.Username,
	})
	c.ChatHistory += fmt.Sprintf("You sent: {%s}\n", fileName)
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
		c.ChatHistory += fmt.Sprintf("%s: %s\n", data.From, data.Message)
	case PrepareFile:
		// When someone is about to send a file
		fmt.Printf("%s: SENT %s \n", data.From, data.Message)
		c.ChatHistory += fmt.Sprintf("%s: SENT {%s}\n", data.From, data.Message)
		// Message should have the filename that's going to be received
		c.fileName = data.Message
	case FileStream:
		// file sent by another user
		fmt.Printf(">>> Downloaded in path {%s}\n", c.fileName)
		c.writeMessageToFile(data)
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
func (c *Client) writeMessageToFile(data *Data) {
	// Write to file
	err := ioutil.WriteFile(DOWNLOAD_FILE_PATH+c.fileName, []byte(data.Message), 0644)

	if err != nil {
		fmt.Println(err)
	}
}
