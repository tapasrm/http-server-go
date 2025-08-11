package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	router := NewRouter()
	router.Add("GET", "/", rootHandler)
	router.Add("GET", "/echo/{message}", echoHandler)

	for {
		conn, err1 := l.Accept()
		if err1 != nil {
			fmt.Println("Error accepting connection: ", err1.Error())
			os.Exit(1)
		}
		go handleConnection(conn, router)
	}
}

func handleConnection(conn net.Conn, router *Router) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

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
		}
		return
	}

	response := router.Find(&req)
	_, err = conn.Write([]byte(response.Serialize()))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}

func rootHandler(req *Request) *Response {
	return &Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": "13"},
		Body:       "Hello, World!",
	}
}

func echoHandler(req *Request) *Response {
	message := req.Params["message"]

	return &Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain", "Content-Length": strconv.Itoa(len(message))},
		Body:       message,
	}
}