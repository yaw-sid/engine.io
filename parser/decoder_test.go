/*
	Decoder tests
*/

package parser

import (
	"testing"

	"github.com/yaw-sid/engineio/frame"
)

var decPacketTestCases = []struct {
	encString       interface{}
	binarySupported bool
	resp            frame.Packet
}{
	{
		encString:       "2probe",
		binarySupported: false,
		resp:            frame.Packet{Type: "ping", Data: "probe"},
	},
	{
		encString:       "b4AQIDBA==",
		binarySupported: false,
		resp:            frame.Packet{Type: "message", Data: []byte{01, 02, 03, 04}},
	},
}

var decPayloadTestCases = []struct {
	pl              string
	binarySupported bool
	resp            frame.Payload
}{
	{
		pl:              "6:4hello2:4$",
		binarySupported: false,
		resp: frame.Payload{
			{Type: "message", Data: "hello"},
			{Type: "message", Data: "$"},
		},
	},
	{
		pl:              "2:4$10:b4AQIDBA==",
		binarySupported: false,
		resp: frame.Payload{
			{Type: "message", Data: "$"},
			{Type: "message", Data: []byte{01, 02, 03, 04}},
		},
	},
}

func TestDecodePacket(t *testing.T) {
	for _, tc := range decPacketTestCases {
		pk, err := DecodePacket(tc.encString, tc.binarySupported)
		if err != nil {
			t.Fatalf("Failed to decode packet: %s", err.Error())
		}
		if pk.Type != tc.resp.Type {
			t.Fatalf("Type inconsistency: %s - %s", pk.Type, tc.resp.Type)
		}
		switch tc.resp.Data.(type) {
		case string:
			if pk.Data != tc.resp.Data {
				t.Fatalf("Data inconsistency: %s - %s", pk.Data, tc.resp.Data)
			}
		case []byte:
			if len(pk.Data.([]byte)) != len(tc.resp.Data.([]byte)) {
				t.Fatalf("Data inconsistency: %d - %d", pk.Data, tc.resp.Data)
			}
		}
		t.Logf("Decoded packet: %v", *pk)
	}
}

func TestDecodePayload(t *testing.T) {
	for _, tc := range decPayloadTestCases {
		pl, err := DecodePayload(tc.pl, tc.binarySupported)
		if err != nil {
			t.Fatalf("Failed to decode payload: %s", err.Error())
		}
		if len(pl) != len(tc.resp) {
			t.Fatalf("Payload inconsistency: %d - %d", len(pl), len(tc.resp))
		}
		for i, pk := range pl {
			if pk.Type != tc.resp[i].Type {
				t.Fatalf("Type inconsistency: %s - %s", pk.Type, tc.resp[i].Type)
			}
			switch tc.resp[i].Data.(type) {
			case string:
				if pk.Data != tc.resp[i].Data {
					t.Fatalf("Data inconsistency: %s - %s", pk.Data, tc.resp[i].Data)
				}
			case []byte:
				if len(pk.Data.([]byte)) != len(tc.resp[i].Data.([]byte)) {
					t.Fatalf("Data inconsistency: %d - %d", pk.Data, tc.resp[i].Data)
				}
			}
		}
	}
}

func BenchmarkDecodePacket(b *testing.B) {
	for _, tc := range decPacketTestCases {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			DecodePacket(tc.encString, tc.binarySupported)
		}
	}
}

func BenchmarkDecodePayload(b *testing.B) {
	for _, tc := range decPayloadTestCases {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			DecodePayload(tc.pl, tc.binarySupported)
		}
	}
}
