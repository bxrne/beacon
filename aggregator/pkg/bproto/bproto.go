package bproto

import (
	"bytes"
	"fmt"
)

const (
	StartByte = 0x02
	EndByte   = 0x03
)

// ParseResponse parses the response according to the byte-aligned protocol.
func ParseResponse(response []byte) (string, error) {
	// Find the end of the HTTP headers
	headersEnd := bytes.Index(response, []byte("\r\n\r\n"))
	var payload []byte

	if headersEnd != -1 {
		// Skip the headers
		payloadStart := headersEnd + 4
		if len(response) <= payloadStart {
			return "", fmt.Errorf("response too short after headers")
		}
		payload = response[payloadStart:]
	} else {
		// No headers found, assume the response is the payload
		payload = response
	}

	if payload[0] != StartByte {
		return "", fmt.Errorf("invalid start byte, got %x instead of %x", payload[0], StartByte)
	}

	length := int(payload[1])

	// Check if we have enough bytes for the full message
	if len(payload) < length+3 {
		return "", fmt.Errorf("response too short for declared length")
	}

	if payload[length+2] != EndByte {
		return "", fmt.Errorf("invalid end byte, got %x instead of %x", payload[length+2], EndByte)
	}

	actualPayload := payload[2 : 2+length]
	return string(actualPayload), nil
}
