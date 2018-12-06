/*
	Package main consists of server and client model

   					The Working
		----------------------------------------

		The server will listen to a spwcific port

		In order to run the server type
		+-------------------------------------------------+
		|	<EXEC_CMD> mode server <PORT_NUM>             |
	+-------------------------------------------------+
*/

package main

import (
	"fmt"
	"os"
	"strconv"
)

/**
 * The default server room
 */
var DEFAULT_SERVER_ROOM = "Server787800p"

/**
 * The HANDLER object handle the entire rooms for chat
 * - It is created when server/client mode is used
 */
var HANDLER = new(Handler)

/**
 * The main function
 * starts by reading the execution commnads
 * app <commad> <param>
 */
func main() {
	args := os.Args

	HANDLER.init()

	if len(args) == 1 {
		fmt.Println("Hey need support!!")
	} else if len(args) == 2 {
		fmt.Println("Ok next time enter command parameter as well as: \n \t app <commad> <param>")
	} else {
		execute(args)
	}

}

/**
 * Execute the instruction as per the arguments provided
 */
func execute(a []string) {
	if a[1] == "mode" { //evaluate the mode section
		executeMode(a[2:])
	} else {
		fmt.Printf("ERR!! command not supported")
	}
}

/**
 * Execute the mode commnad
 */
func executeMode(param []string) {
	fmt.Printf("execute mode command : %s \n", param)

	if param[0] == "server" {

		if len(param) <= 1 {
			fmt.Println("Enter port information")
			return
		}
		PORT, err := strconv.Atoi(param[1])

		if err != nil {
			fmt.Println("Error port num!!!")
			return
		}

		server := Server{}
		server.start(PORT)
		return

	}

	if param[0] == "client" {

		if len(param) <= 1 {
			fmt.Println("Enter IP information \n\t <IP> <PORT> <CLIENT_ID>")
			return
		}

		if len(param) <= 2 {
			fmt.Printf("Enter port information : \n\t%s <PORT> <CLIENT_ID> \n", param[1])
			return
		}

		if len(param) <= 3 {
			fmt.Printf("Enter clietn ID : \n\t%s %s <CLIENT_ID> \n", param[1], param[2])
			return
		}

		var IP, clientID = param[1], param[3]
		PORT, err := strconv.Atoi(param[2])

		if err != nil {
			fmt.Println("Error valid port num!!!")
			return
		}

		client := Client{id: clientID}
		client.connect(IP, PORT)
		return

	}

}
