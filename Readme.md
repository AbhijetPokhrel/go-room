# A simple room based chat in go 

This is a smple room based chat system in go lang. The clients will chat in a specific room. In order to run the system we need to first start the server in a your hosted envronment. The server then listens for clients that will chat in a room.

## Step 1 : Build the project
```
go build
```

## Step 2 : Start the server
```
./build_name mode server <PORT_NUM>
```

## Step 3 : Join server from unique client id and room name
```
./build_name mode client <SERVER_IP> <PORT_NUM> <UNIQUE_CLIENT_ID> <ROOM_NAME>
```
