/*
	Package Main
	The handler that will  handle all the connction and maps
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

/**
 * The handler struct handls the rooms
 */
type Handler struct {
	rooms   map[string]Room      //diffrent Room presetn in the server
	sockets map[string]*net.Conn //diffrest sockets of clients map
}

type Room struct {
	clients []string
}

func (handler *Handler) init() {
	//init handler here
}

func (handler *Handler) handleClient(conn *net.Conn) {

	fmt.Printf("Adding new client\n")

	//read the init message
	for {
		msgByte := make([]byte, 1024)

		len, err := (*conn).Read(msgByte)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot add the client\n")
			break
		}

		if len <= 0 {
			fmt.Println("Cannot add the client %d\n", len)
			break
		}

		fmt.Printf("message : `%s`\n", string(msgByte))
		//check message string here
		// for now just check length
		message := new(Message)
		err = json.Unmarshal(msgByte, &message)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot add the client\n")
			break
		}

		if message.ClientID == "" && message.MSG_TYPE != MsgType["CONTROL_INIT"] {
			fmt.Println("Error message sent for init")
			break
		}

		//every thing is ok so new lets add client to the handler
		client := Client{
			id:   message.ClientID,
			conn: conn,
		}
		//add client to handler
		handler.addClient(&client)

		//start to listen message from the client
		go handler.listenMessage(&client)

	}
}

/**
 * This function will add client to the default handler room of server
 */
func (handler *Handler) addClient(client *Client) {

	handler.sockets[(*client).id] = (*client).conn

}

func (handler Handler) clientDisconnected(clientId string) {
	fmt.Printf("`%s` client disconnected \n")
}

/**
 * Now listen for message
 *
 * Called from the server end
 */
func (handler *Handler) listenMessage(client *Client) {
	fmt.Println("Listening for message from the client\n")
}

func (handler *Handler) initServer(client *Client) {

	var msgByte []byte

	for {
		msgByte = make([]byte, 1024)
		if _, err := (*(*client).conn).Read(msgByte); err == io.EOF {
			handler.clientDisconnected((*client).id)
		}
	}
}
