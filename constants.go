/*
	This is where all the constants are present
*/
package main

// MaxWaitTime  is the time for a single string message to arrive upon
// if the  message does not completely arrive upon this time , the handler will throw an error
var MaxWaitTime int64 = 3000 // 3s

// StrMsgBufferSize is the buffer sizee fo the string message
var StrMsgBufferSize = 800

// +------------------------------### Msg section ###---------------------------------------------+

// MsgSep is the separator that separates each of the messages
// for string message it marks the end of the message
var MsgSep = []byte(">>>====!@#\r\nRRt%^**456--I Love Nepal <<<=========\r\nrty456$7")

var keyValSep = []byte{'>', '-', '+'}

var elemSep = []byte("^\r\n%8r^")

// controlMsg is the control command indicating messsage
// the message can be either control or normal
// control message include operations like initialize, end, restore..etc
var controlMsg = 0xF0

// statusMsg represents the status info
// the status message may be ping messag
// also it may be pong message
// thte user typing is also one satatus
var statusMsg = 0x70

// streamMsg represnt the straming data
// the streaMSg may be file, call, video etc...
var streamMsg = 0x30

// normalMsg is the normal messaging request like string message or file message
var normalMsg = 0x00

// + ------------------------ Diffrent Control messages ------------------------------------+

// ctrlInitMsg is the message type for initializing the connection
var ctrlInitMsg = controlMsg | (0x01) // 1111 0001

// cntrlEndMsg is the message type for termination the connection
var ctrlEndMsg = controlMsg | (0x02) // 1111 0010

// +------------------------  Diffrent Normal messages -------------------------------------+

// norStrMsg is the normal stirng message
var norStrMsg = normalMsg | (0x01) // 0000 0001

// norFileMsg is the normal file message
var norFileMsg = normalMsg | (0x02) // 0000 0010

// norIntroMsg is the normal introduction message
var norIntroMsg = normalMsg | (0x03) // 0000 0011

// norDouYouKnwMsg is the normal message to query other mate if they knw someone
var norDouYouKnwMsg = normalMsg | (0x04) // 0000 0100

// norDouYouKnwMsg is the normal message to query other mate if they knw someone
var norFileReq = normalMsg | (0x05) // 0000 0101

// +------------------------  Diffrent Status messages -------------------------------------+

var statusTyping = statusMsg | (0x00) // 0111 0000
var statusPing = statusMsg | (0x01)   // 0111 0001

// +------------------------  Diffrent Stream messages -------------------------------------+

var streamFile = streamMsg | (0x00) // 0011 0000

// +------------------------------### Msg section Ends ###---------------------------------------------+

// defaultRoom is the default room for chat in this application
var defautRoom = "ROOM_DEF12"
