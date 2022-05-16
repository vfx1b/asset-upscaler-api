package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

func parseImageAndImageConfig(src multipart.File) (image.Image, *image.Config, error) {
	config, _, err := image.DecodeConfig(src)
	if err != nil {
		return nil, nil, err
	}

	// reset read head to original pos
	_, err = src.Seek(0, 0)
	if err != nil {
		return nil, nil, err
	}

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, nil, err
	}

	return img, &config, nil
}
