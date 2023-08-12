package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var logo = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    \".       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n"

type TCPChatServer struct {
	listener          net.Listener
	clients           map[net.Conn]*client
	messageHistory    []string
	rwmutex           *sync.RWMutex
	mutex             *sync.Mutex
	activeConnections int
	usedNames         map[string]bool
}

type client struct {
	Name   string
	Writer *bufio.Writer
}

func NewServer() *TCPChatServer {
	return &TCPChatServer{
		clients:   make(map[net.Conn]*client),
		usedNames: make(map[string]bool),
		rwmutex:   &sync.RWMutex{},
		mutex:     &sync.Mutex{},
	}
}

func (server *TCPChatServer) Listen(typeConnection string, address string) error {
	dstream, err := net.Listen(typeConnection, address)
	if err == nil {
		server.listener = dstream
	}
	log.Printf("Listening on %v", address)
	return err
}

func (server *TCPChatServer) Close() {
	server.listener.Close()
}

func (server *TCPChatServer) Start() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
			continue
		}
		server.mutex.Lock()
		if server.activeConnections >= 10 {
			conn.Close()
			log.Println("Rejected connection due to too many active connections.")
		} else {
			server.activeConnections++
			go server.handleRequest(conn)
		}
		server.mutex.Unlock()
	}
}

func (server *TCPChatServer) broadcast(message string) {
	server.rwmutex.RLock()
	defer server.rwmutex.RUnlock()
	for _, client := range server.clients {
		client.Writer.WriteString(message + "\n")
		client.Writer.Flush()
	}
}

func (server *TCPChatServer) accept(conn net.Conn) *client {
	writer := bufio.NewWriter(conn)
	client := &client{Writer: writer}
	for {
		conn.Write([]byte("Enter your name: "))
		name, _ := bufio.NewReader(conn).ReadString('\n')
		name = strings.TrimSpace(name)
		if name != "" {
			if !server.isNameUsed(name) {
				client.Name = name
				server.markNameAsUsed(name)
				break
			} else {
				conn.Write([]byte("Please, provide a unique name\n"))
			}
		} else {
			conn.Write([]byte("Please, provide non-empty name\n"))
		}
	}
	return client
}

func (server *TCPChatServer) remove(conn net.Conn) {
	server.rwmutex.Lock()
	defer server.rwmutex.Unlock()
	delete(server.clients, conn)
}

func (server *TCPChatServer) uploadHistory(client *client) {
	server.rwmutex.RLock()
	defer server.rwmutex.RUnlock()
	for _, msg := range server.messageHistory {
		client.Writer.WriteString(msg + "\n")
		client.Writer.Flush()
	}
}

func (server *TCPChatServer) welcome(conn net.Conn) {
	conn.Write([]byte(logo))
}

func (server *TCPChatServer) handleRequest(conn net.Conn) {
	defer conn.Close()
	server.welcome(conn)
	client := server.accept(conn)
	server.broadcast(fmt.Sprintf("[Server] %s joined the chat", client.Name))

	server.addClient(conn, client)

	server.uploadHistory(client)

	server.sendMessage(conn)

	server.remove(conn)

	server.broadcast(fmt.Sprintf("[Server] %s left the chat", client.Name))
}

func (server *TCPChatServer) addClient(conn net.Conn, client *client) {
	server.rwmutex.Lock()
	server.clients[conn] = client
	server.rwmutex.Unlock()
}

func (server *TCPChatServer) sendMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message != "" {
			timeStamp := time.Now().Format("2002-07-07 15:04:05")
			fullMessage := fmt.Sprintf("[%s][%s]: %s", timeStamp, server.clients[conn].Name, message)
			server.broadcast(fullMessage)

			server.rwmutex.Lock()
			server.messageHistory = append(server.messageHistory, fullMessage)
			server.rwmutex.Unlock()
		}
	}
}

func (server *TCPChatServer) isNameUsed(name string) bool {
	server.rwmutex.RLock()
	defer server.rwmutex.RUnlock()
	return server.usedNames[name]
}

func (server *TCPChatServer) markNameAsUsed(name string) {
	server.rwmutex.Lock()
	defer server.rwmutex.Unlock()
	server.usedNames[name] = true
}
