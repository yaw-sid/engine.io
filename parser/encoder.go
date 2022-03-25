package parser

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
)

func setPacketType(pk packet) (string, error) {
	// Set packet type
	switch pk.typ {
	case "open":
		return strconv.Itoa(OPEN), nil
	case "close":
		return strconv.Itoa(CLOSE), nil
	case "ping":
		return strconv.Itoa(PING), nil
	case "pong":
		return strconv.Itoa(PONG), nil
	case "message":
		return strconv.Itoa(MESSAGE), nil
	case "upgrade":
		return strconv.Itoa(UPGRADE), nil
	case "noop":
		return strconv.Itoa(NOOP), nil
	default:
		return "", errors.New("invalid packet type")
	}
}

func encodeBuffer(data []byte, packetString string, binarySupported bool) string {
	if binarySupported {
		return string(data)
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	return "b" + packetString + string(dst)
}

// EncodePacket encodes packet into a string. Error is nil if encoding succeeds
func EncodePacket(pk packet, binarySupported bool) (string, error) {
	pkStr, err := setPacketType(pk)
	if err != nil {
		return "", err
	}

	// Encode packet data
	if pk.data != nil {
		switch pk.data.(type) {
		case string:
			pkStr += pk.data.(string)
		case []byte:
			pkStr = encodeBuffer(pk.data.([]byte), pkStr, binarySupported)
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
