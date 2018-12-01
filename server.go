package main

import (
	"fmt"
	"net"
)

/**
 * The server of the application
 */

type Server struct {
	PORT int
}

func (server *Server) start(PORT int) {
	server.PORT = PORT
	listener, error := net.Listen("tcp", fmt.Sprintf(":%d", server.PORT))

	if error != nil {
		fmt.Println(error)
		return
	}
	fmt.Printf("listeining at PORT %d\n", server.PORT)

	// done insists for server termination
	done := make(chan bool)
	go server.listen(listener, done)

	<-done

	// now the server terminates
	// TODO : Add things to do after server termination
	// may be notify all the clients about it
}

/**
 * It is called as a go routine
 */
func (server *Server) listen(listener net.Listener, done chan bool) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go HANDLER.handleClient(conn)
		//every thing is ok
	}

}
