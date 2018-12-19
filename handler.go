package main

import (
	"bytes"
	"errors"
	"fmt"
	"go_chat/helper"
	"net"
	"sync"
)

// Handler handles the client connections
type Handler struct {
	rooms   map[string]*Room   //diffrent Room is present in the server
	clients map[string]*Client //diffrest sockets of clients map
	mutex   *sync.Mutex        //for add and remove client operation
}

// Room is the room where client chat
// A single client can be in multitple rooms
type Room struct {
	clients []*Client
}

// init initializest the handler
func (handler *Handler) init() {
	//init handler here
	handler.rooms = make(map[string]*Room)
	handler.clients = make(map[string]*Client)
	handler.mutex = &sync.Mutex{}
}

// addClient  will add client to the default handler room of server
func (handler *Handler) addClient(client *Client) {

	handler.mutex.Lock()
	fmt.Printf("client added \n")
	// add the client to the map
	handler.clients[(*client).id] = client
	handler.mutex.Unlock()

}

// removeClient removes the client
func (handler *Handler) removeClient(client *Client) {

	handler.mutex.Lock()
	fmt.Printf("removing client( %s )\n", client.id)
	// terminate the client
	client.terminate()
	// delete the client from the map
	delete(handler.clients, client.id)
	handler.mutex.Unlock()

}

// read Reads from the socket
// The main logic here is that we read until we find the msg separator
// If in the read process the time exceeds the MAX_WAIT_TIME an error is thrown
// We will also setDeadline here
func (handler *Handler) read(conn *net.Conn) ([]byte, error) {

	// total is the total writes in bytes
	var total = 0
	// read is the total read in single Read operation
	var read = 0
	// err is the error variable
	var err error
	// msgBytes it the byte on which a read operation writes the data
	var msgByte []byte
	// nowMillis is the current time in milliseconds
	nowMillis := helper.NowAsUnixMilli()
	// buf is the buffer where we append all the read bytes
	var buf bytes.Buffer

	//first setup read deadline

	/**
	 * Loop intil the MSG_SEP is found
	 */
	for {

		// check if the max wait time has exceeded
		if (helper.NowAsUnixMilli() - nowMillis) > MaxWaitTime {
			return nil, errors.New("Read Max execution time exceeded")
		}
		// initialize the message byte to the buffer size
		msgByte = make([]byte, StrMsgBufferSize)
		// read from the client sock
		read, err = (*conn).Read(msgByte)
		// if there is any error return it
		if err != nil {
			return nil, err
		}
		// add total read
		total += read
		// write to the byte buffer
		buf.Write(msgByte[0:read])

		if read > len(MsgSep) {
			if bytes.Equal(MsgSep, msgByte[read-len(MsgSep):read]) {
				total = total - len(MsgSep)
				return buf.Bytes()[0:total], nil
			}
		}

	}

}

// write Writes payload to the socket of the client
func (handler *Handler) write(payload []byte, conn *net.Conn) error {

	// total is the total write in bytes
	var total = 0
	// wlen is the total write in a single Write operation
	var wlen = 0
	// err is the error varaiable
	var err error
	// nowMillis is the current time in milliseconds
	nowMillis := helper.NowAsUnixMilli()

	/**
	 * Loop until all the payload are writeen
	 */
	for {

		// compare if the max write time has exceeded
		if (helper.NowAsUnixMilli() - nowMillis) > MaxWaitTime {
			return errors.New("Write Max execution time exceeded")
		}
		// write to the socket
		wlen, err = (*conn).Write(payload[total:len(payload)])
		if err != nil {
			return errors.New("Cannot write to client")
		}
		// add the total read
		total += wlen
		// check if all the payloa is written
		if total >= len(payload) {
			return nil
		}

	}
}
