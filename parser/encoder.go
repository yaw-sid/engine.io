package parser

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/yaw-sid/engineio/frame"
)

func setPacketType(pk frame.Packet) (string, error) {
	// Set packet type
	switch pk.Type {
	case "open":
		return strconv.Itoa(frame.OPEN), nil
	case "close":
		return strconv.Itoa(frame.CLOSE), nil
	case "ping":
		return strconv.Itoa(frame.PING), nil
	case "pong":
		return strconv.Itoa(frame.PONG), nil
	case "message":
		return strconv.Itoa(frame.MESSAGE), nil
	case "upgrade":
		return strconv.Itoa(frame.UPGRADE), nil
	case "noop":
		return strconv.Itoa(frame.NOOP), nil
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
func EncodePacket(pk frame.Packet, binarySupported bool) (string, error) {
	pkStr, err := setPacketType(pk)
	if err != nil {
		return "", err
	}

	// Encode packet data
	if pk.Data != nil {
		switch pk.Data.(type) {
		case string:
			pkStr += pk.Data.(string)
		case []byte:
			pkStr = encodeBuffer(pk.Data.([]byte), pkStr, binarySupported)
		case map[string]interface{}:
			j, err := json.Marshal(pk.Data.(map[string]interface{}))
			if err != nil {
				return "", err
			}
			pkStr += string(j)
		}
	}

	return pkStr, nil
}

// EncodePayload encodes payload into a string. Error is nil if encoding succeeds
func EncodePayload(p frame.Payload, binarySupported bool) (string, error) {
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
