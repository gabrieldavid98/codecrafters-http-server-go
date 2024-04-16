package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

var dir string

func main() {
	flag.StringVar(&dir, "directory", "", "--directory <directory>")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var (
		req     = readConnToBuffer(conn)
		httpReq = newHttpReq(req)
	)

	responder(conn, httpReq)
}

func readConnToBuffer(conn net.Conn) *bytes.Buffer {
	b := make([]byte, 1024)
	_, err := conn.Read(b)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}

	return bytes.NewBuffer(b)
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
	case strings.HasPrefix(req.path, "/files/"):
		response = respondFiles(req)
	}

	if _, err := conn.Write(response); err != nil {
		fmt.Println("Error writing into connection: ", err.Error())
		os.Exit(1)
	}
}

func respondEcho(req *httpRequest) []byte {
	param, _ := strings.CutPrefix(req.path, "/echo/")
	return okText(param)
}

func respondUserAgent(req *httpRequest) []byte {
	userAgent, ok := req.headers["User-Agent"]
	if !ok {
		return InternalServerErrorResponse
	}

	return okText(userAgent)
}

func respondFiles(req *httpRequest) []byte {
	if req.method == "POST" {
		return respondFilesPost(req)
	}

	return respondFilesGet(req)
}

func respondFilesGet(req *httpRequest) []byte {
	if len(dir) == 0 {
		return NotFoundResponse
	}

	filename, ok := strings.CutPrefix(req.path, "/files/")
	if !ok {
		return InternalServerErrorResponse
	}

	filePath := path.Join(dir, filename)
	f, err := os.Open(filePath)
	if err != nil {
		return NotFoundResponse
	}
	defer f.Close()

	var b bytes.Buffer
	if _, err = io.Copy(&b, f); err != nil {
		return InternalServerErrorResponse
	}

	return okOctetStream(b.Bytes())
}

func respondFilesPost(req *httpRequest) []byte {
	if len(dir) == 0 {
		return NotFoundResponse
	}

	filename, ok := strings.CutPrefix(req.path, "/files/")
	if !ok {
		return InternalServerErrorResponse
	}

	filePath := path.Join(dir, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return NotFoundResponse
	}
	defer f.Close()

	if _, err = io.Copy(f, req.body); err != nil {
		return InternalServerErrorResponse
	}

	return CreatedResponse
}
