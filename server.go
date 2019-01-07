package main

import (
	"fmt"
	"net"
)

/**
 * The server of the application
 */

type Server struct {
	PORT int
}

func (server *Server) start(PORT int) {

	server.PORT = PORT
	listener, error := net.Listen("tcp", fmt.Sprintf(":%d", server.PORT))

	if error != nil {
		fmt.Println(error)
		return
	}
	fmt.Printf("listeining at PORT %d\n", server.PORT)

	// done insists for server termination
	done := make(chan bool)
	go server.listen(listener, done)

	<-done

	// now the server terminates
	// TODO : Add things to do after server termination
	// may be notify all the clients about it
}

// listen listens for the incomming client connections
func (server *Server) listen(listener net.Listener, done chan bool) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go server.handleClient(&conn)
		//every thing is ok
	}

}

// handleClient handles the client
// It listens for clients payload
// Determines what kind of payload is it
// Finally replies as per the payload type
// Payloa is simply the MsgType (Refer constants for MsgType)
func (server *Server) handleClient(conn *net.Conn) {

	fmt.Printf("Adding new client\n")

	client := Client{
		conn: conn,
	}

	//read the init message
	//the init message will not need more than 100 bytes
	var err error
	var msg []byte
	msg, err = handler.read(&client)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot read from the client\n")
		return
	}

	//decode the message
	message, err := _messsageDecode(&msg)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot read from the client\n")
		return
	}

	// if there is no clientId its not the valid message
	if message.ClientID == "" {
		fmt.Println("Error message sent for init")
		return
	}

	err = server.handleMessage(message, &client)

	if err != nil {
		fmt.Println(err)
		return
	}

}

// startListeningFromClient listens from the client
// its loaded from the groutine and will notify the clients state (dead,alive or else)
// it consists of a channel that will mark the termination of the client listening the channel will be the error channek
// once the go routine is terminated the clients go routine will also terminate
func (server *Server) startListeningFromClient(client *Client) error {

	var err error
	var msg []byte

	for {

		msg, err = handler.read(client)

		if err != nil {
			return err
		}

		if err != nil {
			handler.removeClient(client)
		}

		//decode the message

		message, err := _messsageDecode(&msg)

		if err != nil {

			fmt.Printf("\n\n ## msg ========================================= >||||\n  %s \n |||<<<< =========================================\n\n", msg)

			return err
		}

		err = server.handleMessage(message, client)

		if err != nil {
			return err
		}

	}
}

func (server *Server) handleMessage(message *Message, client *Client) error {

	switch message.MsgType {

	case ctrlInitMsg:
		// every thing is ok so new lets add client to the handler
		client.id = message.ClientID
		// now send the client ok message
		message := Message{
			Msg: []byte{'O', 'k'},
		}
		// send ok message
		err := handler.write(message.controlInit(), (*client).conn)
		if err != nil {
			return nil
		}
		// finally add client to handler
		handler.addClient(client)
		// now listen from the client
		err = server.startListeningFromClient(client)
		if err != nil {
			fmt.Println("error 11")
			fmt.Println(err)
			handler.removeClient(client)
			return err
		}
	case norStrMsg:
		// normal string message
		server.processNewMessage(message)

	case statusTyping:
		// user is typing
		server.processNewMessage(message)
	case streamFile:
		// file stream message
		server.processNewMessage(message)

	}

	return nil
}

// processNewMessage will handle the new string message from the client
// first it will check if the room for which the message has arrived is available
// if there is room it will broadcast the message to the room
// if there is no room then it will create the new room
func (server *Server) processNewMessage(message *Message) {

	if room, available := handler.rooms[message.RoomID]; available {
		room.broadCast(message)
	} else { // so there is no such room, lets create one
		handler.createRoom(message.RoomID, message.ClientID)
	}
}
