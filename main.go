package main

// The demo chat applicaiton fro go lang
//
// Abhijet Pokhrel
// (c) 2018

import (
	"fmt"
	"os"
	"strconv"
)

/**
 * The HANDLER object handle the entire rooms for chat
 * - It is created when server/client mode is used
 */
var HANDLER Handler

/**
 * The main function
 * starts by reading the execution commnads
 * app <commad> <param>
 */
func main() {
	args := os.Args
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

	HANDLER = Handler{}

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
			fmt.Println("Enter IP information")
			return
		}

		if len(param) <= 2 {
			fmt.Println("Enter port information")
			return
		}

		IP := param[1]
		PORT, err := strconv.Atoi(param[2])

		if err != nil {
			fmt.Println("Error valid port num!!!")
			return
		}

		client := Client{}
		client.connect(IP, PORT)
		return

	}

}
