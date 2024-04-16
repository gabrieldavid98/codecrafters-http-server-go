package main

import (
	"bytes"
	"fmt"
)

const (
	HttpVersion = "HTTP/1.1"
)

var (
	OkResponse                  = []byte(HttpVersion + " 200 OK\r\n\r\n")
	CreatedResponse             = []byte(HttpVersion + " 201 Created\r\n\r\n")
	NotFoundResponse            = []byte(HttpVersion + " 404 Not Found\r\n\r\n")
	InternalServerErrorResponse = []byte(HttpVersion + " 500 Internal Server Error\r\n\r\n")
)

type httpResponse struct {
	status  int
	headers headers
	body    []byte
}

func (h *httpResponse) withContentLengthHeader() {
	h.headers.addFromKeyValue(
		"Content-Length", fmt.Sprintf("%d", len(h.body)),
	)
}

func (h *httpResponse) withContentTypeHeader(contentType string) {
	h.headers.addFromKeyValue(
		"Content-Type", contentType,
	)
}

func (h *httpResponse) Bytes() []byte {
	var sb bytes.Buffer

	switch h.status {
	case 200:
		sb.WriteString(HttpVersion + " 200 OK\r\n")
	case 404:
		sb.WriteString(HttpVersion + " 404 Not Found\r\n")
	}

	sb.WriteString(h.headers.String())
	sb.WriteString("\r\n")
	sb.Write(h.body)

	return sb.Bytes()
}

func okText(body string) []byte {
	response := &httpResponse{
		status:  200,
		body:    []byte(body),
		headers: newHeaders(),
	}

	response.withContentTypeHeader("text/plain")
	response.withContentLengthHeader()

	return response.Bytes()
}

func okOctetStream(content []byte) []byte {
	response := &httpResponse{
		status:  200,
		body:    content,
		headers: newHeaders(),
	}

	response.withContentTypeHeader("application/octet-stream")
	response.withContentLengthHeader()

	return response.Bytes()
}
