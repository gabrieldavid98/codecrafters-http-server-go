package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

type httpRequest struct {
	method      string
	path        string
	httpVersion string
	headers     headers
	body        *bytes.Buffer
}

func newHttpReq(req *bytes.Buffer) *httpRequest {
	return parseReq(req)
}

func parseReq(req *bytes.Buffer) *httpRequest {
	var (
		httpReq = &httpRequest{
			headers: newHeaders(),
		}
		sc = bufio.NewScanner(req)
	)
	sc.Split(bufio.ScanLines)

	sc.Scan()
	startLine := sc.Text()
	httpReq.parseStartLine(startLine)

	httpReq.parseHeaders(sc)
	httpReq.parseBody(sc)

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

func (h *httpRequest) parseBody(s *bufio.Scanner) {
	contentLength, ok := h.headers["Content-Length"]
	if !ok {
		return
	}

	n, _ := strconv.ParseInt(contentLength, 10, 64)

	if s.Scan() {
		r := bytes.NewReader(s.Bytes())
		h.body = bytes.NewBuffer(nil)
		io.CopyN(h.body, r, n)
	}
}
