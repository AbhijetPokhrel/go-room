/**
 * The message wrapper
 */

package main

type Message struct {
}

var MSG_TYPE = map[string]int{
	"CONTROL": 1,
	"STR_MSG": 2,
}

var CNTRL_MSG_TYPE = map[string]int{
	"INIT": 1,
	"END":  2,
}
