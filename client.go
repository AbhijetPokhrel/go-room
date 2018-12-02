package main

import (
	"encoding/json"
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
	err = client.sendMessage("I am client")

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
func (client *Client) sendMessage(message string) error {

	msgObj := map[string]string{
		"msg":      message,
		"clientId": client.ID,
	}
	jsonObj, err := json.Marshal(msgObj)

	if err != nil {
		return err
	}
	_, err1 := client.conn.Write(jsonObj)

	if err1 != nil {
		return err1
	}
	return nil
}
