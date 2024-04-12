package main

import (
	"fmt"
	"strings"
)

const (
	HttpVersion      = "HTTP/1.1"
	OkResponse       = HttpVersion + " 200 OK\r\n\r\n"
	NotFoundResponse = HttpVersion + " 404 Not Found\r\n\r\n"
)

type httpResponse struct {
	status  int
	headers []header
	body    string
}

type header struct {
	name  string
	value string
}

func (h *httpResponse) withContentLengthHeader() {
	header := header{
		"Content-Length", fmt.Sprintf("%d", len(h.body)),
	}

	h.headers = append(h.headers, header)
}

func (h *httpResponse) withContentTypeHeader(contentType string) {
	header := header{
		"Content-Type", contentType,
	}

	h.headers = append(h.headers, header)
}

func (h *httpResponse) String() string {
	var sb strings.Builder

	switch h.status {
	case 200:
		sb.WriteString(HttpVersion + " 200 OK\r\n")
	case 404:
		sb.WriteString(HttpVersion + " 404 Not Found\r\n")
	}

	for _, header := range h.headers {
		sb.WriteString(header.String())
	}

	sb.WriteString("\r\n")
	sb.WriteString(h.body)

	return sb.String()
}

func (h header) String() string {
	return fmt.Sprintf("%s: %s\r\n", h.name, h.value)
}
