package frame

const (
	// OPEN is sent from the server when a new transport is opened(recheck).
	OPEN = iota
	// CLOSE is sent to request the close of this transport but does not
	// shutdown the connection itself.
	CLOSE
	// PING is sent by the client. Server should answer with a pong packet
	// containing the same data.
	PING
	// PONG is sent by the server to respond to ping packets.
	PONG
	// MESSAGE is the actual message. Client and server should call their callbacks
	// with the data.
	MESSAGE
	// UPGRADE is sent by the client after a test to switch transport is done and
	// requests the server to flush its cache on the old transport and switch to
	// the new transport.
	UPGRADE
	// NOOP is used primarily to force a poll cycle when an incoming websocket
	// connection is received.
	NOOP
)

type Packet struct {
	Type string
	Data interface{}
}
