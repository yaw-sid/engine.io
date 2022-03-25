package parser

import "testing"

var decPacketTestCases = []struct {
	encString       interface{}
	binarySupported bool
	resp            packet
}{
	{
		encString:       "2probe",
		binarySupported: false,
		resp:            packet{typ: "ping", data: "probe"},
	},
	{
		encString:       "b4AQIDBA==",
		binarySupported: false,
		resp:            packet{typ: "message", data: []byte{01, 02, 03, 04}},
	},
}

var decPayloadTestCases = []struct {
	pl              string
	binarySupported bool
	resp            payload
}{
	{
		pl:              "6:4hello2:4$",
		binarySupported: false,
		resp: payload{
			{typ: "message", data: "hello"},
			{typ: "message", data: "$"},
		},
	},
	{
		pl:              "2:4$10:b4AQIDBA==",
		binarySupported: false,
		resp: payload{
			{typ: "message", data: "$"},
			{typ: "message", data: []byte{01, 02, 03, 04}},
		},
	},
}

func TestDecodePacket(t *testing.T) {
	for _, tc := range decPacketTestCases {
		pk, err := DecodePacket(tc.encString, tc.binarySupported)
		if err != nil {
			t.Fatalf("Failed to decode packet: %s", err.Error())
		}
		if pk.typ != tc.resp.typ {
			t.Fatalf("Type inconsistency: %s - %s", pk.typ, tc.resp.typ)
		}
		switch tc.resp.data.(type) {
		case string:
			if pk.data != tc.resp.data {
				t.Fatalf("Data inconsistency: %s - %s", pk.data, tc.resp.data)
			}
		case []byte:
			if len(pk.data.([]byte)) != len(tc.resp.data.([]byte)) {
				t.Fatalf("Data inconsistency: %d - %d", pk.data, tc.resp.data)
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
			if pk.typ != tc.resp[i].typ {
				t.Fatalf("Type inconsistency: %s - %s", pk.typ, tc.resp[i].typ)
			}
			switch tc.resp[i].data.(type) {
			case string:
				if pk.data != tc.resp[i].data {
					t.Fatalf("Data inconsistency: %s - %s", pk.data, tc.resp[i].data)
				}
			case []byte:
				if len(pk.data.([]byte)) != len(tc.resp[i].data.([]byte)) {
					t.Fatalf("Data inconsistency: %d - %d", pk.data, tc.resp[i].data)
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
