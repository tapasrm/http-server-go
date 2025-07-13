package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func (r *Response) SetHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
}

func (r *Response) SetBody(body string) {
	r.Body = body
}

func (r *Response) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) Serialize() string {
	// Start with status line
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, http.StatusText(r.StatusCode))

	// Add headers
	var headerLines strings.Builder
	for key, values := range r.Headers {
		for _, value := range values {
			headerLines.WriteString(fmt.Sprintf("%s: %s\r\n", key, string(value)))
		}
	}

	// End headers with CRLF and add body
	return statusLine + headerLines.String() + "\r\n" + r.Body
}
