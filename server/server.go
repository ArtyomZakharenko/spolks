package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Connection established with %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		fmt.Printf("Received command from %s: %s\n", conn.RemoteAddr(), input)
		parts := strings.SplitN(input, " ", 2)
		command := strings.ToUpper(parts[0])
		var response string

		switch command {
		case "ECHO":
			if len(parts) > 1 {
				response = parts[1]
			} else {
				response = ""
			}
		case "TIME":
			response = time.Now().Format("2006-01-02 15:04:05")
		case "CLOSE", "EXIT", "QUIT":
			response = "Closing connection. Goodbye!"
			conn.Write([]byte(response + "\n"))
			return
		default:
			response = "Unknown command."
		}

		conn.Write([]byte(response + "\n"))
	}

	fmt.Printf("Connection closed with %s\n", conn.RemoteAddr())
}

func startServer(host string, port string) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s:%s\n", host, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

func main() {
	host := "127.0.0.1"
	port := "8080"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	startServer(host, port)
}
