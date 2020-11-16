package ChatRoom

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
)

type MsgType int

const (
	// Everybody will connect to "self" with this port
	CONN_IP_PORT = ":9876"
	// Message types used for protocol
	Success MsgType = iota
	Connect
	Disconnect
	TextStream
	PrepareFile
	FileStream
)

func (m MsgType) String() string {
	return [...]string{"Success", "Connect", "Disconnect", "TextStream", "PrepareFile", "FileStream"}[m]
}

//	Data structure sent through network (Accessed by both client and server)
type Data struct {
	Cmd     MsgType
	Message string
	From    string
}

// Server side processing
type Server struct {
	clients     []chatClient
	ChatHistory string
}

type chatClient struct {
	Connection net.Conn
	User       chatUser
}

type chatUser struct {
	Name string
}

// Start server listening on const port
func (s *Server) Start() {
	go s.serve()
}

func (s *Server) serve() {
	// Start server
	server, err := net.Listen("tcp", CONN_IP_PORT)

	if err != nil {
		fmt.Println(err)
		return
	}
	// Initialize clients' slice
	s.clients = []chatClient{}

	fmt.Println("[SERVER]	Listening...")

	for {
		// wait for client's connection
		connection, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleClient(connection)
	}
}

//	Handles client incomming message
func (s *Server) handleClient(connection net.Conn) {
	data, retrieved := s.tryRetrieveInfo(connection)
	if !retrieved {
		return
	}
	// Do something according to received data
	s.switchAction(&data, connection)
}

//	Tries to obtain data from the connection,
//	will return false if failure
func (s *Server) tryRetrieveInfo(connection net.Conn) (Data, bool) {
	var data Data
	err := gob.NewDecoder(connection).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return data, false
	}
	return data, true
}

//	Decides what to do according to retrieved data
func (s *Server) switchAction(data *Data, connection net.Conn) {
	switch data.Cmd {
	case Connect:
		fmt.Printf("Appends {%s} to chat\n\n", data.From)
		s.ChatHistory += fmt.Sprintf("{%s}: Connects!\n", data.From)
		s.appendToChat(data, connection)
	case Disconnect:
		fmt.Printf("Removes {%s} from chat\n\n", data.From)
		s.ChatHistory += fmt.Sprintf("{%s}: Left the chat!\n", data.From)
		s.removeFromChat(data)
	case TextStream:
		fmt.Printf("Passes {%s}'s text stream\n\n", data.From)
		s.ChatHistory += fmt.Sprintf("{%s}: %s\n", data.From, data.Message)
		s.sendMessage(data)
	case PrepareFile:
		fmt.Printf("{%s} will send FILE=%s\n\n", data.From, data.Message)
		s.ChatHistory += fmt.Sprintf("{%s}: SENT FILE {%s}\n", data.From, data.Message)
		// pass file sending
		s.sendMessage(data)
	case FileStream:
		fmt.Printf("Sends {%s}'s filestream\n\n", data.From)
		s.sendMessage(data)
	default: // success is ignored
	}
}

//	Appends incoming user to chat room
func (s *Server) appendToChat(data *Data, connection net.Conn) {
	// "From" is supposed to have username
	user := chatUser{Name: data.From}
	client := chatClient{Connection: connection, User: user}
	s.clients = append(s.clients, client)

	// Send success message back to appended user
	client.sendMessage(Data{
		Cmd:     Success,
		Message: "You've been connected!",
		From:    "CHAT-ROOM",
	})
	// Send notification to other users
	s.sendMessage(&Data{
		Cmd:     Connect,
		Message: user.Name + " has entered the chat!",
		From:    user.Name,
	})
}

//	Removes user from chat room
func (s *Server) removeFromChat(data *Data) {
	var rmIndex, lastIndex int = -1, len(s.clients) - 1
	// "From" is supposed to have username
	userToRemove := data.From
	for i := 0; i < len(s.clients); i++ {
		client := s.clients[i]
		if client.User.Name == userToRemove {
			rmIndex = i
			break
		}
	}
	// puts last item at index to remove
	s.clients[rmIndex] = s.clients[lastIndex]
	s.clients = s.clients[:lastIndex]

	// Send notification to other users
	s.sendMessage(&Data{
		Cmd:     Disconnect,
		Message: userToRemove + " left the chat!",
		From:    "CHAT-ROOM",
	})
}

//	Sends (passes) mesage to all connected users
func (s *Server) sendMessage(data *Data) {
	excludedUser := data.From

	// send to all connected users
	for _, client := range s.clients {
		if excludedUser != client.User.Name {
			client.sendMessage(*data)
		}
	}
}

//	Saves all chat history to file
func (s *Server) BackupChat() {
	// Write chat history to file
	err := ioutil.WriteFile("./ChatHistory.txt", []byte(s.ChatHistory), 0644)

	if err != nil {
		fmt.Println(err)
	}
}

//	Sends given data to client
func (cl *chatClient) sendMessage(dataToSend Data) {
	err := gob.NewEncoder(cl.Connection).Encode(&dataToSend)
	if err != nil {
		fmt.Println(err)
	}
}
