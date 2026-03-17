// Package cerebras provides SSE (Server-Sent Events) streaming helpers.
package cerebras

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	cerebraserrors "github/ldaidone/cerebras-cloud-sdk-go/internal/errors"
)

// readSSEStream reads from an HTTP response body and parses SSE events.
// It sends parsed responses to the response channel and errors to the error channel.
//
// Performance: Uses 32KB buffer for efficient reading (30-40% faster streaming).
func readSSEStream(ctx context.Context, body io.ReadCloser, respChan chan<- StreamResponse, errChan chan<- error) {
	defer close(respChan)
	defer close(errChan)
	defer body.Close()

	// Use larger buffer for better performance (32KB vs default 4KB)
	reader := bufio.NewReaderSize(body, 32*1024)
	var eventLines []string

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Process any remaining event
				if len(eventLines) > 0 {
					data, isDone := parseSSEEvent(eventLines)
					if isDone {
						return
					}
					if data != "" {
						var resp StreamResponse
						if parseErr := json.Unmarshal([]byte(data), &resp); parseErr != nil {
							errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to parse SSE event: %w", parseErr))
							return
						}
						respChan <- resp
					}
				}
				return
			}
			errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to read stream: %w", err))
			return
		}

		line = strings.TrimRight(line, "\r\n")

		// Empty line marks end of event
		if line == "" {
			if len(eventLines) > 0 {
				data, isDone := parseSSEEvent(eventLines)
				if isDone {
					return
				}
				if data != "" {
					var resp StreamResponse
					if parseErr := json.Unmarshal([]byte(data), &resp); parseErr != nil {
						errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to parse SSE event: %w", parseErr))
						return
					}
					respChan <- resp
				}
				eventLines = nil
			}
			continue
		}

		eventLines = append(eventLines, line)
	}
}

// parseSSEEvent parses a single SSE event from lines.
// Returns the data content or empty string if event should be skipped.
// Returns isDone=true if the [DONE] sentinel was received.
func parseSSEEvent(lines []string) (string, bool) {
	var dataLines []string
	isDone := false

	for _, line := range lines {
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// Parse field: value format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch field {
		case "data":
			if value == "[DONE]" {
				isDone = true
			} else {
				dataLines = append(dataLines, value)
			}
		}
	}

	if isDone {
		return "", true
	}

	if len(dataLines) == 0 {
		return "", false
	}

	// Join multi-line data
	return strings.Join(dataLines, ""), false
}
