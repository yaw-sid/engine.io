/*
	Encoder tests
*/

package parser

import (
	"testing"

	"github.com/yaw-sid/engineio/frame"
)

var packetTestCases = []struct {
	pk              frame.Packet
	binarySupported bool
	response        string
}{
	{
		frame.Packet{Type: "ping", Data: "probe"},
		false,
		"2probe",
	},
	{
		frame.Packet{Type: "message", Data: "$"},
		false,
		"4$",
	},
}

var payloadTestCases = []struct {
	plyd            frame.Payload
	binarySupported bool
	response        string
}{
	{
		frame.Payload{
			{Type: "message", Data: "hello"},
			{Type: "message", Data: "$"},
		},
		false,
		"6:4hello2:4$",
	},
	{
		frame.Payload{
			{Type: "message", Data: "$"},
			{Type: "message", Data: []byte{01, 02, 03, 04}},
		},
		false,
		"2:4$10:b4AQIDBA==",
	},
	{
		frame.Payload{
			{Type: "open", Data: map[string]interface{}{
				"sid":          "lv_VI97HAXpY6yYWAAAC",
				"upgrades":     []string{"websocket"},
				"pingInterval": 25000,
				"pingTimeout":  5000,
			}},
		},
		false,
		"96:0{\"pingInterval\":25000,\"pingTimeout\":5000,\"sid\":\"lv_VI97HAXpY6yYWAAAC\",\"upgrades\":[\"websocket\"]}",
	},
}

func TestEncodePacket(t *testing.T) {
	for _, tc := range packetTestCases {
		pkEnc, err := EncodePacket(tc.pk, tc.binarySupported)
		if err != nil {
			t.Errorf("Failed to encode packet: %s", err.Error())
		}
		if pkEnc != tc.response {
			t.Errorf("Response inconsistency: %s - %s", pkEnc, tc.response)
		}
		t.Logf("Encoded packet: %s", pkEnc)
	}
}

func TestEncodePayload(t *testing.T) {
	for _, tc := range payloadTestCases {
		payloadStr, err := EncodePayload(tc.plyd, tc.binarySupported)
		if err != nil {
			t.Errorf("Failed to encode payload: %s", err.Error())
		}
		if payloadStr != tc.response {
			t.Errorf("Response inconsistency: %s - %s", payloadStr, tc.response)
		}
		t.Logf("Encoded packet: %s", payloadStr)
	}
}

func BenchmarkEncodePacket(b *testing.B) {
	for _, tc := range packetTestCases {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			EncodePacket(tc.pk, tc.binarySupported)
		}
	}
}

func BenchmarkEncodePayload(b *testing.B) {
	for _, tc := range payloadTestCases {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			EncodePayload(tc.plyd, tc.binarySupported)
		}
	}
}
