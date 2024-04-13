package main

import (
	"bufio"
	"strings"
)

type httpRequest struct {
	method      string
	path        string
	httpVersion string
	headers     headers
}

func newHttpReq(reqStr string) *httpRequest {
	return parseReq(reqStr)
}

func parseReq(reqStr string) *httpRequest {
	var (
		httpReq = &httpRequest{
			headers: newHeaders(),
		}
		sc = bufio.NewScanner(strings.NewReader(reqStr))
	)
	sc.Split(bufio.ScanLines)

	sc.Scan()
	startLine := sc.Text()
	httpReq.parseStartLine(startLine)

	httpReq.parseHeaders(sc)

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

func (h *httpRequest) parseHeaders(s *bufio.Scanner) {
	for s.Scan() {
		t := s.Text()

		if len(t) == 0 {
			break
		}

		h.headers.addFromString(t)
	}
}
