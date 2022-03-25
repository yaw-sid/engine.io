/*
	Encoder tests
*/

package parser

import (
	"testing"
)

var packetTestCases = []struct {
	pk              packet
	binarySupported bool
	response        string
}{
	{
		packet{typ: "ping", data: "probe"},
		false,
		"2probe",
	},
	{
		packet{typ: "message", data: "$"},
		false,
		"4$",
	},
}

var payloadTestCases = []struct {
	plyd            payload
	binarySupported bool
	response        string
}{
	{
		payload{
			{typ: "message", data: "hello"},
			{typ: "message", data: "$"},
		},
		false,
		"6:4hello2:4$",
	},
	{
		payload{
			{typ: "message", data: "$"},
			{typ: "message", data: []byte{01, 02, 03, 04}},
		},
		false,
		"2:4$10:b4AQIDBA==",
	},
	{
		payload{
			{typ: "open", data: map[string]interface{}{
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
