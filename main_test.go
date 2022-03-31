package main

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/png"
	"testing"
)

func TestEncodeImageToBase64PNG(t *testing.T) {
	expects := "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"
	expectsBytes, _ := base64.StdEncoding.DecodeString(expects)
	reader := bytes.NewReader(expectsBytes)

	img, _, _ := image.Decode(reader)
	imgb64, _ := encodeImageToBase64PNG(img)

	if expects != imgb64 {
		t.Fatal("Broken image encoding")
	}
}
