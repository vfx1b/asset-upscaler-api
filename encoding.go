package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
)

type ImageEncoding string

const (
	EncoderPNG  = "PNG"
	EncoderJPEG = "JPEG"
)

type ImageStringEncoder interface {
	encodeImage(i image.Image) (string, error)
}

type PngImageStringEncoder struct {
}

func (e *PngImageStringEncoder) encodeImage(image image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

type JpegImageStringEncoder struct {
}

func (e *JpegImageStringEncoder) encodeImage(image image.Image) (string, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, image, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
