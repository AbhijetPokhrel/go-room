package main

import (
	"bytes"
	"errors"
	"fmt"
	"go-room/helper"
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
	id      string
	mutex   *sync.Mutex
	clients []*Client
}

// init initializes the room
func (room *Room) init() {
	room.clients = []*Client{}
	room.mutex = &sync.Mutex{}
}

// broadCast will broadcast a string meesage to the room
func (room *Room) broadCast(message *Message) {

	fmt.Printf("broadcasting to room %s(%d)\n", room.id, len(room.clients))
	found := false
	for _, client := range room.clients {
		client.addMsgQueue(message)
		if !found {
			if client.id == message.ClientID {
				found = true
			}
		}
	}

	if !found {
		room.addClient(handler.clients[message.ClientID])
	}

}

// addClient adds client to the room
// also addes the room to the clients room array
func (room *Room) addClient(client *Client) {

	room.mutex.Lock()
	defer room.mutex.Unlock()

	fmt.Printf("adding client : %s to room : %s \n", client.id, room.id)
	room.clients = append(room.clients, client)
	client.rooms = append(client.rooms, room)

}

// removeClient removes the client from the room
func (room *Room) removeClient(clientID string) {

	room.mutex.Lock()
	fmt.Printf("remove client : %s from room : %s \n", clientID, room.id)

	for i, v := range room.clients {
		if v.id == clientID { // client found so delete it and break the loop
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			break
		}
	}

	room.mutex.Unlock()
}

// clientsNum returns the total number of clients in a room
func (room *Room) clientsNum() int {

	room.mutex.Lock()
	defer room.mutex.Unlock()
	return len(room.clients)

}

// init initializest the handler
func (handler *Handler) init() {
	//init handler here
	handler.rooms = make(map[string]*Room)
	handler.clients = make(map[string]*Client)
	handler.mutex = &sync.Mutex{}
}

// createRoom creates a new room in the server
func (handler *Handler) createRoom(roomID string, clientID string) {

	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	fmt.Printf("Creating a room : %s\n", roomID)
	room := Room{id: roomID}
	room.init()
	room.addClient(handler.clients[clientID])
	handler.rooms[roomID] = &room

}

func (handler *Handler) removeRoom(roomID string) {

	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	delete(handler.rooms, roomID)

}

// addClient  will add client to the default handler room of server
func (handler *Handler) addClient(client *Client) {

	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	fmt.Printf("client added \n")
	// add the client to the map
	client.init()
	handler.clients[(*client).id] = client

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
func (handler *Handler) read(client *Client) ([]byte, error) {

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

		// first of all read form the client buffer
		// the client buffer has some pending messages
		// after reading from the client buffer, read from the socket
		if len(client.buffer) > 0 {

			// initialize message byte
			msgByte = client.buffer
			// read from the client buffer
			read = len(client.buffer)

		} else { // read from the socket here

			// initialize the message byte to the buffer size
			msgByte = make([]byte, StrMsgBufferSize)
			// read from the client sock
			read, err = (*(*client).conn).Read(msgByte)
			// if there is any error return it
			if err != nil {
				return nil, err
			}

		}

		// add total read
		total += read
		// write to the byte buffer
		buf.Write(msgByte[0:read])

		if read > len(MsgSep) {

			b := buf.Bytes()
			index := bytes.Index(b, MsgSep)

			if index != -1 { // -1 markes the index is not found

				if index == 0 {
					return []byte{}, nil
				}

				// if the total read is less than index then it markes that some
				// other messages are also pending from this client
				// we need to listen to that as well
				// hence we need to add these unread messages to the buffer of the client
				if index+len(MsgSep) < total { // add the remaining to the client buffer
					client.buffer = b[index+len(MsgSep) : total]
				} else { // no no more pending message to read
					client.buffer = []byte{}
				}

				return b[0:index], nil
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
