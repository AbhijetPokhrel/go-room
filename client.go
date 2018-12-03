package main

import (
	"fmt"
	"net"
)

/**
 * The client section of the application
 */

type Client struct {
	id   string    // the client ID
	conn *net.Conn // the server sock for client side and client sock for server side
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

	err = client.sendMessage(message.controlInit())

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to the server \n")

}

/**
 * Called when message is received
 */
func (client *Client) onMessageReceived() {

}

/**
 * Send Message
 */
func (client *Client) sendMessage(message []byte) error {

	_, err := (*client.conn).Write(message)

	if err != nil {
		return err
	}
	return nil
}
