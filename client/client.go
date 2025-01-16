package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func startClient(host string, port string) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to server %s:%s\n", host, port)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command (ECHO, TIME, CLOSE/EXIT/QUIT): ")
		if !scanner.Scan() {
			break
		}
		command := scanner.Text()
		if strings.TrimSpace(command) == "" {
			continue
		}

		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Printf("Error sending command: %v\n", err)
			break
		}

		responseScanner := bufio.NewScanner(conn)
		if responseScanner.Scan() {
			fmt.Printf("Server response: %s\n", responseScanner.Text())
		}

		if strings.ToUpper(command) == "CLOSE" || strings.ToUpper(command) == "EXIT" || strings.ToUpper(command) == "QUIT" {
			break
		}
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

	startClient(host, port)
}
