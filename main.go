package main

import (
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

func index(w http.ResponseWriter, r *http.Request) {
	it, _ := template.ParseFiles("index.html")
	err := it.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}
}

func result(w http.ResponseWriter, r *http.Request) {
	// possible memory attack vector
	err := r.ParseMultipartForm(5 << 20)

	scaleFactorString := r.FormValue("scaleFactor")
	scaleFactor, err := strconv.Atoi(scaleFactorString)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	// possible memory attack vector
	if scaleFactor < 1 || scaleFactor > 20 {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	original, config, err := retrieveImageAndMeta(r)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	upscaled := imaging.Resize(original, config.Width*scaleFactor, config.Height*scaleFactor, imaging.NearestNeighbor)

	originalB64, err := encodeImageToBase64PNG(original)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}
	upscaledB64, err := encodeImageToBase64PNG(upscaled)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	templateValues := map[string]string{
		"original": originalB64,
		"upscaled": upscaledB64,
	}

	rt, _ := template.ParseFiles("result.html")
	err = rt.Execute(w, templateValues)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}
}

func retrieveImageAndMeta(r *http.Request) (image.Image, *image.Config, error) {
	src, _, err := r.FormFile("assetFile")
	if err != nil {
		return nil, nil, err
	}

	config, _, err := image.DecodeConfig(src)
	if err != nil {
		return nil, nil, err
	}

	src.Seek(0, 0) // reset read head to original pos

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, nil, err
	}

	return img, &config, nil
}

func encodeImageToBase64PNG(image image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
