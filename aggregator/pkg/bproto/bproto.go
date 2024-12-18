package bproto

import (
	"fmt"
)

// Note: The byte-aligned protocol is not used when requesting /metric.
// It is retained here in case it's needed for other endpoints.

// ParseResponse parses the response according to the byte-aligned protocol.
func ParseResponse(response []byte) (string, error) {
	if len(response) < 3 {
		return "", fmt.Errorf("response too short")
	}
	if response[0] != 0x02 {
		return "", fmt.Errorf("invalid start byte")
	}
	length := int(response[1])
	if length != len(response)-3 {
		return "", fmt.Errorf("length byte does not match payload length")
	}
	payload := response[2 : 2+length]
	if response[2+length] != 0x03 {
		return "", fmt.Errorf("invalid end byte")
	}
	return string(payload), nil
}
