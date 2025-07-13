package main

import (
	"errors"
	"strings"
)

// Some security-related constants
const (
	MAX_HEADERS    = 100
	MAX_URI_LENGTH = 8192
)

type Request struct {
	Method   string              `json:"method"`
	Path     string              `json:"path"` // Changed from map[string]string for simplicity and correctness
	Protocol string              `json:"protocol"`
	Query    map[string][]string `json:"query"`  // Changed to handle multiple values for the same key (prevents HPP)
	Headers  map[string][]string `json:"headers"`// Changed to handle multiple values for the same key
	Params   map[string]string   `json:"params"`
	Host     string              `json:"host"`
	Body     string              `json:"body"`
}

func (req *Request) ParseRequest(request string) error {
	// Find the end of the headers section
	headerBodySplit := strings.SplitN(request, "\r\n\r\n", 2)
	headerBlock := headerBodySplit[0]

	var body string
	if len(headerBodySplit) > 1 {
		body = headerBodySplit[1]
	}
	req.Body = body

	lines := strings.Split(headerBlock, "\r\n")
	if len(lines) < 1 {
		return errors.New("invalid request format: no request line")
	}

	// --- Request-Line Parsing ---
	reqLine := strings.Split(lines[0], " ")
	if len(reqLine) != 3 {
		return errors.New("invalid request line format")
	}

	// Add a basic check for URI length to prevent some DoS attacks
	if len(reqLine[1]) > MAX_URI_LENGTH {
		return errors.New("URI too long")
	}

	req.Method = reqLine[0]
	req.Protocol = reqLine[2]

	// --- Path and Query Parsing ---
	fullPath := reqLine[1]
	pathParts := strings.SplitN(fullPath, "?", 2)
	req.Path = pathParts[0]

	if len(pathParts) > 1 {
		queryParts := strings.Split(pathParts[1], "&")
		req.Query = make(map[string][]string)
		for _, part := range queryParts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				// Append values for the same key to prevent HTTP Parameter Pollution
				req.Query[kv[0]] = append(req.Query[kv[0]], kv[1])
			} else if len(kv) == 1 {
				// Handle keys without values, e.g., ?foo&bar
				req.Query[kv[0]] = append(req.Query[kv[0]], "")
			}
		}
	}

	// --- Header Parsing ---
	headerLines := lines[1:]
	if len(headerLines) > MAX_HEADERS {
		return errors.New("too many headers")
	}

	req.Headers = make(map[string][]string)
	for _, line := range headerLines {
		if line == "" {
			continue // Should not happen with the split, but good practice
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			key := strings.TrimSpace(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			req.Headers[key] = append(req.Headers[key], value)
		}
	}

	// --- Host Header Validation ---
	// The http.Header map canonicalizes keys, but we are not using it.
	// So we need to check for "Host" case-insensitively or rely on clients sending it as "Host".
	// For now, we'll assume the key is "Host" as in the original code.
	// A robust implementation would iterate and check with strings.EqualFold.
	hostHeaders, ok := req.Headers["Host"]
	if !ok || len(hostHeaders) == 0 {
		return errors.New("Host header is required")
	}
	req.Host = hostHeaders[0] // Use the first host header

	return nil
}
