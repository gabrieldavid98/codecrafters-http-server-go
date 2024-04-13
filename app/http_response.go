package main

import (
	"fmt"
	"strings"
)

const (
	HttpVersion                 = "HTTP/1.1"
	OkResponse                  = HttpVersion + " 200 OK\r\n\r\n"
	NotFoundResponse            = HttpVersion + " 404 Not Found\r\n\r\n"
	InternalServerErrorResponse = HttpVersion + " 500 Internal Server Error\r\n\r\n"
)

type httpResponse struct {
	status  int
	headers headers
	body    string
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

func (h *httpResponse) String() string {
	var sb strings.Builder

	switch h.status {
	case 200:
		sb.WriteString(HttpVersion + " 200 OK\r\n")
	case 404:
		sb.WriteString(HttpVersion + " 404 Not Found\r\n")
	}

	sb.WriteString(h.headers.String())
	sb.WriteString("\r\n")
	sb.WriteString(h.body)

	return sb.String()
}

func okText(body string) string {
	response := &httpResponse{
		status:  200,
		body:    body,
		headers: newHeaders(),
	}

	response.withContentTypeHeader("text/plain")
	response.withContentLengthHeader()

	return response.String()
}
