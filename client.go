package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net"
)

// Client is the struct containign client infos
type Client struct {
	id           string     // the client ID
	messageQueue *list.List // the queues of string message
	conn         *net.Conn  // the server sock for client side and client sock for server side
}

/**
 * Connect to a server
 */
func (client *Client) connect(IP string, PORT int) {

	fmt.Printf("Connecting to server at PORT %d\n", PORT)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", IP, PORT))

	if err != nil {
		fmt.Println(err)
		return
	}

	client.conn = &conn

	message := Message{
		ClientID: client.id,
	}

	err = client.sendMessage(&message)

	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan bool)
	go client.listenForServer(done)
	<-done

}

// listenForServer listens for the incommin message from the server
func (client *Client) listenForServer(done chan bool) {

	var msg []byte
	var err error
	var message = new(Message)
	for {

		msg, err = handler.read(client.conn)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot read from the server\n")
			return
		}

		err = json.Unmarshal(msg, &message)

		if err != nil {
			fmt.Printf("Didn't get proper response form server")
			fmt.Println(err)
		}

		// every thing is ok
		if message.MsgType == ctrlInitMsg {
			fmt.Printf("Success: %s \n")
			// initiate message queue
			client.messageQueue = list.New()
		}

		fmt.Printf(" >  : %s", message.Msg)
	}

}

// addMsgQueue adds the message to the clients to be sent by the server
func (client *Client) addMsgQueue(message *Message) {
	client.messageQueue.PushBack(message)
}

// sendMessage sends the message to the server
func (client *Client) sendMessage(message *Message) error {

	return handler.write(message.controlInit(), client.conn)

}
