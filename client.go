package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/Workiva/go-datastructures/queue"
)

// Client is the struct containign client infos
// id is the id of the client
// messageQueue is the queue of messages that the client is waiting to be sent
// conn is the client socket conn
// rooms is the array of rooms the client is associated with
// mutex is the mutex to make sure the message sends are thread safe
// isBusy marks if the cleint is busy sending the message from the message queue
// buffer is the messages that are still left to be read
type Client struct {
	id           string       // the client ID
	messageQueue *queue.Queue // the queues of string message
	conn         *net.Conn    // the server sock for client side and client sock for server side
	rooms        []*Room      // the rooms that its associated with
	mutex        *sync.Mutex  // mutex to make sure the message sends are thread safe
	isBusy       bool         // is client busy sending messages
	buffer       []byte       // the byte of message to be read
}

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Hybrid elements used by both server and client                               |
// +----------------------------------------------------------------------------------------------------------------+

// init initailizes the client
// if called from the server end it marks that the new client is valid and can be initialized
// if called from the client end it marks the server has accepted the client
func (client *Client) init() {
	client.messageQueue = queue.New(10)
	client.rooms = []*Room{}
	client.mutex = &sync.Mutex{}
}

// sendMessage sends the byte message to the server/client
func (client *Client) sendMessage(message *Message) error {
	return handler.write(message.generateByte(), client.conn)
}

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Hybrid elements used by both server and client ends                          |
// +----------------------------------------------------------------------------------------------------------------+

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Server side elements                                                         |
// +----------------------------------------------------------------------------------------------------------------+

// ## Server side element are the function and variable that are use on the server end

// addMsgQueue adds the message to the clients queue to be sent by the server
// This part must be thread safe hence we have added the isBusy clause to mark if the client is busy
// popping and pushing the messages
// If the client is not busy then the queue message are sent bu sendQueueMsgs function
func (client *Client) addMsgQueue(message *Message) {
	if client.id == message.ClientID { // this message is sent by us so no need to add to queue
		return
	}
	client.messageQueue.Put(message)
	// if the client is not busy ask to send the message from the message queue
	if !client.isBusy {
		client.sendQueueMsgs()
	}
}

// sendQueueMsgs function sends the messages in the queue
// In this process the client is marked busy
func (client *Client) sendQueueMsgs() {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	// set client busy
	client.isBusy = true
	var message []interface{}
	var err error

	// loop until the message queue is empty
	for !client.messageQueue.Empty() {
		message, err = client.messageQueue.Get(1)
		if err != nil {
			break
		}
		client.sendMessage(message[0].(*Message))
	}

	// finally set the client free
	client.isBusy = false

}

// terminate terminates a client
// this is done when the server cannot write or read from the client
// TODO terminate on terminate command (YET to add)
func (client *Client) terminate() {

	fmt.Printf("terminating the client : %s\n", client.id)
	for _, room := range client.rooms {
		room.removeClient(client.id)
	}

}

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Server side elements Ends                                                    |
// +----------------------------------------------------------------------------------------------------------------+

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Client side elements                                                         |
// +----------------------------------------------------------------------------------------------------------------+

// ## client side elements are the functions and variables that are used on the client ends

// scanner scanns from buffer input from the terminal
var scanner *bufio.Scanner

// connect connects to the specific ip
// In the connecting process a init message is sent with cntrlInitMsg type (see constants.go)
// Once the connection is established a hello message is sent inorder to synchronize rooms i.e
// create room if its not there in server and if there notify all room members that you have joined
// IP is the IP of the server
// PORT is the PORT in which we are connecting
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
		MsgType:  ctrlInitMsg,
	}

	err = client.sendMessage(&message)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.sendHelloMessage()

	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan bool)
	go client.listenForServer(done)
	client.waitForUserInp()

}

// sendHelloMessage sends the hello message to the server
// The main reason to sendHelloMessage is to make sure the room we are trying to chat is there in server
// If the room is not there a room is created and if the room is there the the hello message is sent to
// all the room members
func (client *Client) sendHelloMessage() error {

	fmt.Printf("Sending hello message \n")
	message := Message{
		ClientID: client.id,
		Msg:      []byte("Hi! I am new user"),
		MsgType:  norStrMsg,
		RoomID:   defautRoom,
	}
	return client.sendMessage(&message)
}

// listenForServer listens for the incommin message from the server
func (client *Client) listenForServer(done chan bool) {

	var msg []byte
	var err error
	var message *Message

	// loop until the EOF reacheds
	for {

		msg, err = handler.read(client)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("Cannot read from the server\n")
			return
		}

		message, err = _messsageDecode(&msg)

		if err != nil {
			fmt.Printf("Didn't get proper response form server\n")
			fmt.Println(err)
		}

		// every thing is ok
		if message.MsgType == ctrlInitMsg {

			// initiate message queue
			client.init()

		} else if message.MsgType == norStrMsg {

			// here is the new string message hence print it
			client.printNewMsg(message)

		} else if message.MsgType == statusTyping {

			if bytes.Equal(message.Msg, []byte("true")) {

				fmt.Printf("%s typing.... \n", message.ClientID)

			} else {

				fmt.Printf("%s stopped typing.... \n", message.ClientID)

			}

		} else if message.MsgType == streamFile {

		}

	}

}

// waitForUserInp waits for user input from the terminal
func (client *Client) waitForUserInp() {

	scanner = bufio.NewScanner(os.Stdin)
	var text string

	fmt.Print(" | ------------- Welcome to the system ----------  |\n")
	for text != "q" { // break the loop if text == "q"

		fmt.Printf("\r You : ")

		scanner.Scan()
		text = scanner.Text()

		if text != "q" {
			message := Message{
				ClientID: client.id,
				Msg:      []byte(text),
				MsgType:  norStrMsg,
				RoomID:   defautRoom,
			}
			client.sendMessage(&message)
		}
	}

}

// printNewMsg prints the new message arrived to the client in the terminal
func (client *Client) printNewMsg(message *Message) {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	fmt.Printf("\033[2K\n\r %s : %s\n\r You : ", message.ClientID, message.Msg)

}

// +----------------------------------------------------------------------------------------------------------------+
// |								 # Client side elements Ends                                                    |
// +----------------------------------------------------------------------------------------------------------------+
