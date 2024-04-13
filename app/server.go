package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	handleConn(conn)
}

func handleConn(conn net.Conn) {
	reqStr := readConnToString(conn)
	fmt.Println(reqStr)

	httpReq := newHttpReq(reqStr)
	responder(conn, httpReq)
}

func readConnToString(conn net.Conn) string {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}

	return string(buff)
}

func responder(conn net.Conn, req *httpRequest) {
	response := NotFoundResponse

	switch {
	case req.path == "/":
		response = OkResponse
	case strings.HasPrefix(req.path, "/echo/"):
		response = respondEcho(req)
	case req.path == "/user-agent":
		response = respondUserAgent(req)
	}

	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Println("Error writing into connection: ", err.Error())
		os.Exit(1)
	}
}

func respondEcho(req *httpRequest) string {
	param, _ := strings.CutPrefix(req.path, "/echo/")
	return okText(param)
}

func respondUserAgent(req *httpRequest) string {
	userAgent, ok := req.headers["User-Agent"]
	if !ok {
		return InternalServerErrorResponse
	}

	return okText(userAgent)
}
