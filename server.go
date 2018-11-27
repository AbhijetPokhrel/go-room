package main

import (
	"fmt"
)

/**
 * The server of the application
 */

type Server struct {
	roomId string
}

func (server *Server) start(PORT int) {
	fmt.Printf("start server of roomId : %s at %d", server.roomId, PORT)
}

func (server *Server) listen
