package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/Workiva/go-datastructures/queue"
)

// Client is the struct containign client infos
type Client struct {
	id           string       // the client ID
	messageQueue *queue.Queue // the queues of string message
	conn         *net.Conn    // the server sock for client side and client sock for server side
	rooms        []*Room      // the rooms that its associated with
}

// init initailizes the client
func (client *Client) init() {
	client.messageQueue = queue.New(10)
}

// connect connects to the specific ip
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

	err = client.sendInit(&message)

	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan bool)
	go client.listenForServer(done)
	client.waitForUserInp()

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
			fmt.Printf("Success: \n")
			// initiate message queue
			client.init()
		}

		fmt.Printf(" >  : %s \n", message.Msg)
	}

}

func (client *Client) waitForUserInp() {

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != "q" { // break the loop if text == "q"
		fmt.Print(" > : ")
		scanner.Scan()
		text = scanner.Text()
		if text != "q" {
			message := Message{
				ClientID: client.id,
				Msg:      text,
			}
			client.sendMessage(&message)
		}
	}

}

// addMsgQueue adds the message to the clients to be sent by the server
func (client *Client) addMsgQueue(message *Message) {
	client.messageQueue.Put(message)
}

// terminate terminates a client
func (client *Client) terminate() {
	// TODO
}

// sendMessage sends the message to the server
func (client *Client) sendInit(message *Message) error {

	return handler.write(message.controlInit(), client.conn)

}

func (client *Client) sendMessage(message *Message) error {
	return handler.write(message.str(), client.conn)
}
