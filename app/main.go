package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Route struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for true {
		conn, err1 := l.Accept()
		if err1 != nil {
			fmt.Println("Error accepting connection: ", err1.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buffer[:n]))

	var req Request
	if err := req.ParseRequest(string(buffer[:n])); err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		response := Response{
			StatusCode: 400,
			Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": "15"},
			Body:       "Bad Request",
		}
		_, err = conn.Write([]byte(response.Serialize()))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
		fmt.Println("Response sent to client")
		return
	}
	fmt.Printf("Parsed Request: Method: %s, Path: %s, Protocol: %s, Host: %s\n",
		req.Method, req.Path, req.Protocol, req.Host)

	fmt.Println("Headers:")
	for key, values := range req.Headers {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("Query Parameters:")
	for key, values := range req.Query {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("Body:", req.Body)

	method, path, protocol := req.Method, req.Path, req.Protocol
	fmt.Printf("Method: %s, Path: %s, Protocol: %s\n", method, path, protocol)

	var response Response
	if path == "/" {
		response.SetStatusCode(200)
		response.SetHeader("Content-Type", "text/plain")
		response.SetBody("Hello, World!")
	}

	if strings.HasPrefix(path, "/echo") {
		params := strings.Split(path, "/")
		pathParam := ""
		if len(params) >= 3 {
			pathParam = params[2]
		}
		byteLen := len([]byte(pathParam))
		response.SetStatusCode(200)
		response.SetHeader("Content-Type", "text/plain")
		response.SetHeader("Content-Length", strconv.Itoa(byteLen))
		response.SetBody(pathParam)
	}
	fmt.Println("Response: ", response)
	_, err = conn.Write([]byte(response.Serialize()))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}
	fmt.Println("temp sent to client")
}
