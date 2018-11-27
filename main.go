/**
 * The demo for chat application in go lang for local network
 *
 * Abhijet Pokhrel
 * (c) 2018
 */
package main

import (
	"fmt"
	"os"
)

/**
 * The HANDLER object handle the entire rooms for chat
 * - It is created when server/client mode is used
 */
var HANDLER: Handler,
    PORT:=5050

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

	if(param[0] == "server"){
		server := Server{
			roomId : param[1]
		}

		server.start(PORT)
	}

}
