package main

import (
	"fmt"
	"net"
)

/**
 * The client section of the application
 */

type Client struct {
	ID   string   // the client ID
	conn net.Conn // the server sock for client side and client sock for server side
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

	client.conn = conn
	if client.sendMessage("I am client") {
		fmt.Println("Connected to the server")
	}

}

/**
 * Called when message is received
 */
func (client *Client) onMessageReceived() {

}

/**
 * Send Message
 */
func (client *Client) sendMessage(message string) bool {

	_, err := client.conn.Write([]byte(message))

	if err != nil {
		fmt.Println("Cannot connect to the server")
		return false
	}
	return true
}
