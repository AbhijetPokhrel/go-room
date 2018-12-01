/*
	Package Main
	The handler that will  handle all the connction and maps
*/
package main

import (
	"fmt"
	"net"
	"unicode/utf8"
)

/**
 * The handler struct handls the rooms
 */
type Handler struct {
	rooms   map[string]Room     //diffrent Room presetn in the server
	sockets map[string]net.Conn //diffret clients socket
}

type Room struct {
	clientId []string
}

func (handler *Handler) handleClient(conn net.Conn) {

	fmt.Printf("Adding new client")

	//read the init message
	for {
		message := make([]byte, 1024)

		len, err := conn.Read(message)

		if err != nil {
			fmt.Println("Cannot add the client\n")
			break
		}

		if len <= 0 {
			fmt.Println("Cannot add the client\n")
			break
		}

		msgStr := string(message)
		fmt.Printf("Init Msg : %s\n", msgStr)

		//check message string here
		// for now just check length

		if utf8.RuneCountInString(msgStr) > 0 {
			go handler.listenMessage(Client{conn: conn})
			break
		}

	}
}

/**
 * Now listen for message
 *
 * Called from the server end
 */
func (handler *Handler) listenMessage(client Client) {
	fmt.Println("Listening for message from the client\n")
}

func (handler *Handler) initServer(client Client) {

}
