package main

import (
	"encoding/json"
	"fmt"
)

// Message is the prototype of message payload
// It is formatted in the message byte to send to the server
// ClientID is the id of the sender
// ReceiverID is the id of the receiver client
// RoomID is the id of the room for which it was assigned for
// Msg is the string message
// MsgType is the type of the message that is being sent (For more details on message type see constants)
type Message struct {
	ClientID   string
	ReceiverID string
	RoomID     string
	Msg        string
	MsgType    int
}

// str generates the message byte for string message
func (message *Message) str() []byte {

	message.MsgType = norStrMsg

	return message.generateByte()
}

// controlInit generates the control message for initiating establishing a connection with the server
func (message *Message) controlInit() []byte {

	message.MsgType = ctrlInitMsg

	return message.generateByte()
}

// generateByte generates the necessary byte along with the server
func (message *Message) generateByte() []byte {

	json, err := json.Marshal(*message)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return append(json, MsgSep...)

}
