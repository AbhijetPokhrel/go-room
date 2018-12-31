package main

import (
	"bytes"
	"errors"
	"strconv"
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
	Msg        []byte
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

	return message.encode()

}

func (message *Message) encode() []byte {

	var buf bytes.Buffer

	if message.ClientID != "" {
		_messageEncodeElem("ClientID", message.ClientID, &buf)
	}

	if message.RoomID != "" {
		_messageEncodeElem("RoomID", message.RoomID, &buf)
	}

	if message.ReceiverID != "" {
		_messageEncodeElem("ReceiverID", message.ReceiverID, &buf)
	}

	if message.MsgType != 0 {
		_messageEncodeElem("MsgType", strconv.Itoa(message.MsgType), &buf)
	}

	buf.Write(message.Msg) // write the msg payload
	buf.Write(MsgSep)

	return buf.Bytes()

}

// _messageDecode decodes the message byte to message pointer
// b is the byte of message
func _messsageDecode(b *[]byte) (*Message, error) {

	message := Message{}

	msg := bytes.Split(*b, elemSep)
	// if the length of msg slice is less than the message is invalid
	if len(msg) < 2 {
		return nil, errors.New("Invalid message")
	}

	// elemCount counts the number of elements added to the message like MsgType,Msg etc..
	// the elemCount should be equal to len(msg) after the loop coming
	var elemCount int

	// loop until the last element
	// the last element is the payload
	for index, element := range msg {

		if (index + 1) == len(msg) {

			message.Msg = element
			elemCount++
			break
		}

		elem := bytes.Split(element, keyValSep)

		if len(elem) < 2 {
			return nil, errors.New("Invalid message")
		}

		// find the approprite elem of message
		// if unknown elem is sent then this is an errors
		switch string(elem[0]) {

		case "ClientID":

			message.ClientID = string(elem[1])
			elemCount++

		case "ReceiverID":

			message.ReceiverID = string(elem[1])
			elemCount++

		case "RoomID":

			message.RoomID = string(elem[1])
			elemCount++

		case "MsgType":

			msgType, err := strconv.ParseInt(string(elem[1]), 10, 16)

			if err != nil {
				return nil, err
			}

			message.MsgType = int(msgType)
			elemCount++

		default: // unknown elemetn which is a error
			return nil, errors.New("Invalid message")

		} // switch case ends

	} // for loop ends

	if elemCount != len(msg) {
		return nil, errors.New("Invalid message")
	}

	// Now we have a valid message

	return &message, nil

}

func _messageEncodeElem(key string, val string, buff *bytes.Buffer) {

	buff.WriteString(key)
	buff.Write(keyValSep)
	buff.WriteString(val)
	buff.Write(elemSep)

}
