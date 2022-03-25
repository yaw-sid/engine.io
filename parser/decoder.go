package parser

// DecodePacket decodes a string into a packet. Error is nil if decoding succeeds
func DecodePacket(pkStr string, binarySupported bool) (*packet, error) {
	var pk packet
	// Set packet type
	switch pkStr[0] {
	case OPEN:
		pk.typ = "open"
	case CLOSE:
		pk.typ = "close"
	case PING:
		pk.typ = "ping"
	case PONG:
		pk.typ = "pong"
	case MESSAGE:
		pk.typ = "message"
	case UPGRADE:
		pk.typ = "upgrade"
	case NOOP:
		pk.typ = "noop"
	}

	// Decode packet data
	if len(pkStr) > 1 {
		if pkStr[1] == 'b' {

		}
		pk.data = pkStr[1:]
	}

	return &pk, nil
}
