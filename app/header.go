package main

import (
	"strings"
)

type headers map[string]string

func newHeaders() headers {
	return make(headers)
}

func (h headers) addFromString(s string) {
	p := strings.Split(s, ": ")

	if len(p) == 0 {
		return
	}

	h[p[0]] = p[1]
}

func (h headers) addFromKeyValue(key, value string) {
	h[key] = value
}

func (h headers) String() string {
	var sb strings.Builder

	for k, v := range h {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\r\n")
	}

	return sb.String()
}
