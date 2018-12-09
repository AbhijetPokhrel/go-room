package main

import (
	"encoding/json"
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

		go handleClient(&conn)
		//every thing is ok
	}

}

// handleClient handles the client
// It listens for clients payload
// Determines what kind of payload is it
// Finally replies as per the payload type
// Payloa is simply the MsgType (Refer constants for MsgType)
func handleClient(conn *net.Conn) {

	fmt.Printf("Adding new client\n")

	//read the init message
	//the init message will not need more than 100 bytes
	var err error
	var msg []byte
	msg, err = handler.read(conn)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot read from the client\n")
		return
	}

	fmt.Printf("msg :` %s`", msg)

	//decode the message
	message := new(Message)
	err = json.Unmarshal(msg, &message)

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

	err = handleMessage(message, conn)

	if err != nil {
		fmt.Println(err)
	}

}

func handleMessage(message *Message, conn *net.Conn) error {

	switch message.MsgType {

	case ctrlInitMsg:
		// every thing is ok so new lets add client to the handler
		client := Client{
			id:   message.ClientID,
			conn: conn,
		}
		// now send the client ok message
		message := Message{
			Msg: "OK",
		}
		// send ok message
		err := handler.write(message.controlInit(), conn)
		if err != nil {
			return nil
		}
		// finally add client to handler
		handler.addClient(&client)

	}
	return nil
}
