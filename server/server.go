package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var uploadsDir = "uploads"
var connections sync.Map

type FileTransfer struct {
	FileName string
	Offset   int64
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	clientID := conn.RemoteAddr().String()
	fmt.Printf("Connection established with %s\n", clientID)

	var currentTransfer *FileTransfer
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		fmt.Printf("Received command from %s: %s\n", clientID, input)
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
		case "UPLOAD":
			if len(parts) < 2 {
				response = "Usage: UPLOAD <filename>"
				break
			}
			fileName := parts[1]
			filePath := filepath.Join(uploadsDir, fileName)
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				response = "Error creating file: " + err.Error()
				break
			}
			if currentTransfer != nil {
				file.Seek(currentTransfer.Offset, io.SeekStart)
			}
			conn.Write([]byte("READY\n"))
			written, err := io.Copy(file, conn)
			if err != nil {
				fmt.Printf("Error during upload: %v\n", err)
				currentTransfer = &FileTransfer{FileName: fileName, Offset: written}
				response = "Connection lost during upload. File partially saved."
				break
			}
			file.Close()
			response = "File uploaded successfully."
		case "DOWNLOAD":
			if len(parts) < 2 {
				response = "Usage: DOWNLOAD <filename>"
				break
			}
			fileName := parts[1]
			filePath := filepath.Join(uploadsDir, fileName)
			file, err := os.Open(filePath)
			if err != nil {
				response = "Error opening file: " + err.Error()
				break
			}
			fileInfo, _ := file.Stat()
			conn.Write([]byte("READY " + strconv.FormatInt(fileInfo.Size(), 10) + "\n"))
			_, err = io.Copy(conn, file)
			if err != nil {
				fmt.Printf("Error during download: %v\n", err)
				response = "Connection lost during download."
				break
			}
			file.Close()
			response = "File downloaded successfully."
		case "CLOSE", "EXIT", "QUIT":
			response = "Closing connection. Goodbye!"
			conn.Write([]byte(response + "\n"))
			return
		default:
			response = "Unknown command."
		}

		if response != "" {
			conn.Write([]byte(response + "\n"))
		}
	}

	fmt.Printf("Connection closed with %s\n", clientID)
}

func startServer(host string, port string) {
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

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
