package main

import (
	"bufio"
	"strings"
)

type httpRequest struct {
	method      string
	path        string
	httpVersion string
}

func newHttpReq(reqStr string) *httpRequest {
	return parseReq(reqStr)
}

func parseReq(reqStr string) *httpRequest {
	httpReq := new(httpRequest)

	reqScanner := bufio.NewScanner(strings.NewReader(reqStr))
	reqScanner.Split(bufio.ScanLines)

	reqScanner.Scan()
	startLine := reqScanner.Text()
	httpReq.parseStartLine(startLine)

	return httpReq
}

func (h *httpRequest) parseStartLine(startLine string) {
	s := bufio.NewScanner(strings.NewReader(startLine))
	s.Split(bufio.ScanWords)

	s.Scan()
	h.method = s.Text()

	s.Scan()
	h.path = s.Text()

	s.Scan()
	h.httpVersion = s.Text()
}
