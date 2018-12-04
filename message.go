/**
 * The message wrapper
 */

package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ClientID string
	Msg      string
	MSG_TYPE int
}

var MsgType = map[string]int{
	"CONTROL_INIT": 100,
	"CONTROL_END":  101,
	"STR_MSG":      200,
}

var MSG_SEP = []byte{'%', '/', 'y', '#', '!'}

/**
 * Generate the string message ie normal message
 */
func (message *Message) str(msg string) []byte {
	message.Msg = msg

	message.MSG_TYPE = MsgType["STR_MSG"]

	return generateByte(message)
}

/**
 * Gnereate control message
 */

func (message *Message) controlInit() []byte {

	message.MSG_TYPE = MsgType["CONTROL_INIT"]

	return generateByte(message)
}

func generateByte(message *Message) []byte {

	json, err := json.Marshal(*message)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return append(json, MSG_SEP...)

}
