package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"./gpio"
)

var clients []Client

type Message struct {
	client  Client
	message string
}

type Client struct {
	connection net.Conn
}

func (c Client) Read(ch chan<- Message) {
	bufc := bufio.NewReader(c.connection)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			fmt.Println("Force closing connection - Unable to read from connection: ", err)
			c.connection.Close()
			break
		}
		ch <- Message{client: c, message: line}
	}
}

func (c Client) Write(message Message) {
	_, err := io.WriteString(c.connection, message.message)
	if err != nil {
		fmt.Println("Unable to write to connection: ", err)
		return
	}
}

func handleMessage(messageChannel <-chan Message) {
	for {
		message := <-messageChannel
		fmt.Printf("A message has been receieved: '%s'", message.message)
		if message.message == "up\n" {
			fmt.Println("WE GOT AN UP!")
			gpio.WriteToPin("15")
//			gpio.WriteToPin("18")
			gpio.WriteToPin("17")
		} else if message.message == "down\n" {
			fmt.Println("WE GOT A DOWN!")
		} else if message.message == "left\n" {
			fmt.Println("WE GOT A LEFT!")
		} else if message.message == "right\n" {
			fmt.Println("WE GOT A RIGHT!")
		} else if message.message == "stop\n" {
                        fmt.Println("WE GOT A RIGHT!")			
			gpio.Reset()
		}

		for _, c := range clients {
			if c != message.client {
				go c.Write(message)
			}
		}
	}
}

func handleConnection(connection net.Conn, messageChannel chan<- Message) {
	client := Client{connection: connection}
	clients = append(clients, client)
	go logConnection("Connected", connection)
	go client.Read(messageChannel)
}

func logConnection(event string, connection net.Conn) {
	fmt.Printf("%v - %v.\n", connection.RemoteAddr(), event)
}

func main() {
	gpio.Reset()
	defer gpio.Reset()
	flag.Parse()
	port := flag.Arg(0)

	// Start listening
	ln, err := net.Listen("tcp", ":"+port)

	// handle errors resulting from a failed listen attempt
	if err == nil {
		fmt.Printf("Server listening on %v\n", ln.Addr())
	} else {
		fmt.Println("Unable to start listener: ", err)
		os.Exit(1)
	}

	clients = make([]Client, 0)
	messageChannel := make(chan Message)
	go handleMessage(messageChannel)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept incoming connection: ", err)
			continue
		}

		go handleConnection(conn, messageChannel)
	}
}
