package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var MAX_WAIT_TIME = 3000 // 3s

/**
 * The handler struct handls the rooms
 */
type Handler struct {
	rooms   map[string]*Room     //diffrent Room presetn in the server
	sockets map[string]*net.Conn //diffrest sockets of clients map
}

type Room struct {
	clients []string
}

func (handler *Handler) init() {
	//init handler here
	handler.rooms = make(map[string]*Room)
	handler.sockets = make(map[string]*net.Conn)
}

func (handler *Handler) handleClient(conn *net.Conn) {

	fmt.Printf("Adding new client\n")

	//read the init message
	//the init message will not need more than 100 bytes
	msgByte := make([]byte, 100)
	var err error
	var read int
	var nowMillis = NowAsUnixMilli()
	read, err = handler.read(&nowMillis, &msgByte, conn)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot add the client\n")
		return
	}

	//decode the message
	message := new(Message)
	err = json.Unmarshal(msgByte[0:read], &message)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot add the client\n")
		return
	}

	if message.ClientID == "" && message.MSG_TYPE != MsgType["CONTROL_INIT"] {
		fmt.Println("Error message sent for init")
		return
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

// Read from the socket
// TODO : add timer for read
func (handler *Handler) read(nowMillis *int64, msgByte *[]byte, conn *net.Conn) (int, error) {

	var total = 0
	var read = 0
	var err error

	for {
		if (NowAsUnixMilli() - *nowMillis) > int64(MAX_WAIT_TIME) {
			return 0, errors.New("Max execution time exceeded")
		}

		read, err = (*conn).Read(*msgByte)
		total += read

		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot add the client\n")
			return 0, err
		}

		if read <= 0 {
			fmt.Println("Cannot add the client %d\n", read)
			return 0, err
		}

		if total > len(MSG_SEP) {

			if bytes.Equal(MSG_SEP, (*msgByte)[total-len(MSG_SEP):total]) {
				total = total - len(MSG_SEP)
				return total, nil
			}

		}

	}

}

func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
