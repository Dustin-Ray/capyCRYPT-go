package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

func RunClient() {
	// Create a tls.Config object with the desired tls properties
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	// Connect to the server
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", tlsConfig)
	if err != nil {
		fmt.Printf("Error dialing server: %s\n", err)
		return
	}
	defer conn.Close()

	// Make sure the connection is established by printing out the server's name
	state := conn.ConnectionState()
	fmt.Println("Connected to server:", state.ServerName)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
