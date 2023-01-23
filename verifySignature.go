package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
)

func verifySignature(signature, timestamp, body string) bool {
	var msg bytes.Buffer

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	pubKey, err := hex.DecodeString(PUBLIC_KEY)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	msg.WriteString(timestamp)
	msg.WriteString(body)

	return ed25519.Verify(pubKey, msg.Bytes(), sig)
}
