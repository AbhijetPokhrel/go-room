/*
	Package main consists of server and client model

	+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

	# Working

		This project demonstrates the use of room based chat system.
		- The server will be started first in a specific port
		- This port will then listen for incomming connections
		- The client now will connect to the server specifying the room it want to join
		- If the room is not there then the handler which is run by server instance will create a new room

	# The server

		The server will listen to the specific port

		# In order to run the server , type

		+-------------------------------------------+
		#	<EXEC_CMD> mode server <PORT_NUM>
		+-------------------------------------------+

	# The client

		The client will connect to the server by specifying the host port client id and room name

	+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

*/
package main
