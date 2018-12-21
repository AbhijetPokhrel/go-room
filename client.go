package main

import (
	"bufio"
	"encoding/json"
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
type Client struct {
	id           string       // the client ID
	messageQueue *queue.Queue // the queues of string message
	conn         *net.Conn    // the server sock for client side and client sock for server side
	rooms        []*Room      // the rooms that its associated with
	mutex        *sync.Mutex  // mutex to make sure the message sends are thread safe
	isBusy       bool         // is client busy sending messages
}

// init initailizes the client
func (client *Client) init() {
	client.messageQueue = queue.New(10)
	client.rooms = []*Room{}
	client.mutex = &sync.Mutex{}
}

var scanner *bufio.Scanner

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
		MsgType:  ctrlInitMsg,
	}

	err = client.sendMessage(&message)

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
			fmt.Printf("Cannot read from the server\n")
			return
		}

		err = json.Unmarshal(msg, &message)

		if err != nil {
			fmt.Printf("Didn't get proper response form server\n")
			fmt.Println(err)
		}

		// every thing is ok
		if message.MsgType == ctrlInitMsg {
			// initiate message queue
			client.init()
		} else if message.MsgType == norStrMsg {

			client.printNewMsg(message)
		}

	}

}

func (client *Client) waitForUserInp() {

	scanner = bufio.NewScanner(os.Stdin)
	var text string

	fmt.Print(" | ------------- Welcome to the system ----------  |\n")
	for text != "q" { // break the loop if text == "q"

		fmt.Printf("\rYou : ")

		scanner.Scan()
		text = scanner.Text()

		if text != "q" {
			message := Message{
				ClientID: client.id,
				Msg:      text,
				MsgType:  norStrMsg,
				RoomID:   defautRoom,
			}
			client.sendMessage(&message)
		}
	}

}

func (client *Client) printNewMsg(message *Message) {

	client.mutex.Lock()
	defer client.mutex.Unlock()
	fmt.Printf("\r\n")
	fmt.Printf("\033[A\033[A\r %s : %s \n", message.ClientID, message.Msg)

}

// addMsgQueue adds the message to the clients queue to be sent by the server
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

func (client *Client) sendQueueMsgs() {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.isBusy = true
	fmt.Printf("client ( %s ) busy \n", client.id)
	var message []interface{}
	var err error
	for !client.messageQueue.Empty() {
		message, err = client.messageQueue.Get(1)
		if err != nil {
			break
		}
		client.sendMessage(message[0].(*Message))
	}

	client.isBusy = false
	fmt.Printf("client ( %s ) free \n", client.id)

}

// terminate terminates a client
func (client *Client) terminate() {

	fmt.Printf("terminating the client : %s\n", client.id)
	for _, room := range client.rooms {
		room.removeClient(client.id)
	}

}

// sendMessage sends the byte message to the server/client
func (client *Client) sendMessage(message *Message) error {
	return handler.write(message.generateByte(), client.conn)
}
