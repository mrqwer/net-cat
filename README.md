# Alem01 - Netcat 
## TCPChat - Group Chat Server and Client

TCPChat is a simple group chat implementation using a Server-Client architecture. It allows multiple clients to connect to a central server and exchange messages in a group chat format. The server handles incoming connections, manages client messages, and broadcasts messages to all connected clients.

## Features 
- TCP connection between the server and multiple clients (one-to-many).
- Clients are required to provide a name.
- Control the maximum number of allowed connections (10 connections).
- Clients can send messages to the chat.
- Empty messages are not broadcasted.
- Messages are timestamped with the sender's name and time of sending.
- New clients receive previous chat history upon joining.
- Clients are notified when another client joins or leaves.
- All clients receive messages sent by other clients.
- Clients can disconnect without affecting other clients.
- Default port is 8989, but can be specified as a command line argument.

## Usage
Compile the project:
```bash
go build -o tcpchat main.go
```

Run the server:
```bash
./tcpchat [port]
```
* `port` (optional): The port number to listen on. If not provided, the default port is 8989 is used.

Clients can connect to the server using `nc` command:
```bash
nc localhost [port]
```

Each client is prompted to enter a name. Once connected, clients can start sending messages to the chat.


To exit, simply close the client's terminal.

## Example
```bash
./tcpchat 8888
```
```bash
nc localhost 8888
```

## Dependencies
Go version 1.20.1
Used only the standard library

## References
[Bao Nguyen](https://nqbao.medium.com/writing-a-chat-server-in-go-3b61ccc2a8ed)

[Kazuki Higashiguchi](https://dev.to/hgsgtk/how-go-handles-network-and-system-calls-when-tcp-server-1nbd)


## Authors
[@mrqwer](https://github.com/mrqwer)