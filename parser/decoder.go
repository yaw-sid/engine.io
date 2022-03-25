package parser

import (
	"encoding/base64"
	"errors"
	"strconv"
)

func decodeBase64String(str string) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(str)))
	n, err := base64.StdEncoding.Decode(dst, []byte(str))
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// DecodePacket decodes a string into a packet. Error is nil if decoding succeeds
func DecodePacket(encPk interface{}, binarySupported bool) (*packet, error) {
	switch encPk.(type) {
	case []byte:
		return &packet{typ: "message", data: encPk}, nil
	case string:
		if encPk.(string)[0] == 'b' {
			decodedBase64, err := decodeBase64String(encPk.(string)[2:])
			if err != nil {
				return nil, err
			}
			return &packet{typ: "message", data: decodedBase64}, nil
		}
		// Set packet type
		var pk packet
		switch string(encPk.(string)[0]) {
		case strconv.Itoa(OPEN):
			pk.typ = "open"
		case strconv.Itoa(CLOSE):
			pk.typ = "close"
		case strconv.Itoa(PING):
			pk.typ = "ping"
		case strconv.Itoa(PONG):
			pk.typ = "pong"
		case strconv.Itoa(MESSAGE):
			pk.typ = "message"
		case strconv.Itoa(UPGRADE):
			pk.typ = "upgrade"
		case strconv.Itoa(NOOP):
			pk.typ = "noop"
		}
		if len(encPk.(string)) > 1 {
			pk.data = encPk.(string)[1:]
		}
		return &pk, nil
	}
	return nil, errors.New("invalid packet")
}

func splitPayload(pl string) ([]string, error) {
	var packets []string
	var n int
	var err error
	var length string
	var msg string
	for i := 0; i < len(pl); {
		ch := pl[i]
		if string(ch) != ":" {
			length += string(ch)
			i++
			continue
		}
		n, err = strconv.Atoi(length)
		if err != nil {
			return nil, err
		}
		if length == "" || length != strconv.Itoa(n) {
			return nil, errors.New("parser error - ignoring payload")
		}
		msg = pl[(i + 1):(n + i + 1)]
		if length != strconv.Itoa(len(msg)) {
			return nil, errors.New("parser error - ignoring payload")
		}
		packets = append(packets, msg)
		i += n + 1
		length = ""
	}
	return packets, nil
}

// DecodePayload decodes a string into a payload. Error is nil if decoding succeeds
func DecodePayload(encPayload string, binarySupported bool) (payload, error) {
	encPackets, err := splitPayload(encPayload)
	if err != nil {
		return nil, err
	}

	var pyld payload

	for _, encPk := range encPackets {
		decPk, err := DecodePacket(encPk, binarySupported)
		if err != nil {
			return nil, err
		}
		pyld = append(pyld, *decPk)
	}

	return pyld, nil
}
