package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
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

func checkInputFormatsOrFail(r *ResultHandlerForm, w http.ResponseWriter) {
	// if r.OutputImageFormat == EncoderJPEG || r.OutputImageFormat == EncoderPNG {
	// 	return
	// }

	// http.Error(w, "Internal Error", http.StatusBadRequest)
}

func checkInputImagesOrFail(r *ResultHandlerForm, w http.ResponseWriter) {
	for _, i := range r.Images {
		if !i.Valid {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}
	}
}
