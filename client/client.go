package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
		fmt.Print("Enter command (ECHO, TIME, UPLOAD <file>, DOWNLOAD <file>, CLOSE/EXIT/QUIT): ")
		if !scanner.Scan() {
			break
		}
		command := scanner.Text()
		if strings.TrimSpace(command) == "" {
			continue
		}

		parts := strings.SplitN(command, " ", 2)
		switch strings.ToUpper(parts[0]) {
		case "UPLOAD":
			if len(parts) < 2 {
				fmt.Println("Usage: UPLOAD <filename>")
				continue
			}
			fileName := parts[1]
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				continue
			}
			conn.Write([]byte(command + "\n"))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			if strings.TrimSpace(response) == "READY" {
				_, err := io.Copy(conn, file)
				if err != nil {
					fmt.Printf("Error during file upload: %v\n", err)
				}
				fmt.Println("File upload complete.")
			}
			file.Close()
		case "DOWNLOAD":
			if len(parts) < 2 {
				fmt.Println("Usage: DOWNLOAD <filename>")
				continue
			}
			conn.Write([]byte(command + "\n"))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			if strings.HasPrefix(strings.TrimSpace(response), "READY") {
				fileSize, _ := strconv.ParseInt(strings.TrimSpace(strings.Split(response, " ")[1]), 10, 64)
				file, err := os.Create(parts[1])
				if err != nil {
					fmt.Printf("Error creating file: %v\n", err)
					continue
				}
				io.CopyN(file, conn, fileSize)
				file.Close()
				fmt.Println("File download complete.")
			} else {
				fmt.Println(strings.TrimSpace(response))
			}
		default:
			conn.Write([]byte(command + "\n"))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Printf("Server response: %s\n", strings.TrimSpace(response))
		}

		if strings.ToUpper(parts[0]) == "CLOSE" || strings.ToUpper(parts[0]) == "EXIT" || strings.ToUpper(parts[0]) == "QUIT" {
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
