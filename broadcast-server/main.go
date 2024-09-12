package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
)

var connections = make(map[net.Conn]bool)
var mutex = sync.Mutex{}

func broadcastMessage(msg string, conn net.Conn) {
	for c := range connections {
		if c != conn {
			_, err := c.Write([]byte(msg))
			if err != nil {
				fmt.Println(err)
				fmt.Println("Could not send message to", c.RemoteAddr())
			}
		}
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("New connection from", conn.RemoteAddr())

	defer func() {
		conn.Close()
		mutex.Lock()
		delete(connections, conn)
		mutex.Unlock()
		broadcastMessage(fmt.Sprintf("User %s has left the server\n", conn.RemoteAddr()), conn)
	}()

	mutex.Lock()
	connections[conn] = true
	mutex.Unlock()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		broadcastMessage(msg, conn)
	}
}

func listenServerShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			broadcastMessage("Server is shutting down\n", nil)
			os.Exit(0)
		}
	}()
}

func main() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
		return
	}

	listenServerShutdown()

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}
