package parser

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
)

// EncodePacket encodes packet into a string. Error is nil if encoding succeeds
func EncodePacket(pk packet, binarySupported bool) (string, error) {
	var pkStr string
	// Set packet type
	switch pk.typ {
	case "open":
		pkStr = strconv.Itoa(OPEN)
	case "close":
		pkStr = strconv.Itoa(CLOSE)
	case "ping":
		pkStr = strconv.Itoa(PING)
	case "pong":
		pkStr = strconv.Itoa(PONG)
	case "message":
		pkStr = strconv.Itoa(MESSAGE)
	case "upgrade":
		pkStr = strconv.Itoa(UPGRADE)
	case "noop":
		pkStr = strconv.Itoa(NOOP)
	default:
		return "", errors.New("invalid packet type")
	}

	// Encode packet data
	if pk.data != nil {
		switch pk.data.(type) {
		case string:
			pkStr += pk.data.(string)
		case []byte:
			if binarySupported {
				for _, ch := range pk.data.([]byte) {
					pkStr += strconv.Itoa(int(ch))
				}
				break
			}
			dst := make([]byte, base64.StdEncoding.EncodedLen(len(pk.data.([]byte))))
			base64.StdEncoding.Encode(dst, pk.data.([]byte))
			pkStr = "b" + pkStr + string(dst)
		case map[string]interface{}:
			j, err := json.Marshal(pk.data.(map[string]interface{}))
			if err != nil {
				return "", err
			}
			pkStr += string(j)
		}
	}

	return pkStr, nil
}

// EncodePayload encodes payload into a string. Error is nil if encoding succeeds
func EncodePayload(p payload, binarySupported bool) (string, error) {
	var payloadStr string
	// Encode each packet in the payload
	for _, pk := range p {
		pkStr, err := EncodePacket(pk, binarySupported)
		if err != nil {
			return "", err
		}
		payloadStr += strconv.Itoa(len(pkStr)) + ":" + pkStr
	}

	return payloadStr, nil
}
