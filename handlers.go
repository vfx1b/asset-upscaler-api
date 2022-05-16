package main

import (
	"net/http"
	"text/template"

	"github.com/disintegration/imaging"
	"github.com/ggicci/httpin"
)

func index(w http.ResponseWriter, r *http.Request) {
	it, _ := template.ParseFiles("index.html")
	err := it.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}
}

type ResultHandlerForm struct {
	Images            []httpin.File `in:"form=images"`
	ScaleFactor       int           `in:"form=scale-factor"`
	OutputImageFormat string        `in:"form=output-image-format"`
}

func result(w http.ResponseWriter, r *http.Request) {
	form := r.Context().Value(httpin.Input).(*ResultHandlerForm)

	encoder := imageEncoders[ImageEncoding(form.OutputImageFormat)]
	if encoder == nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	for _, i := range form.Images {
		if !i.Valid {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}
	}

	// possible memory attack vector
	if form.ScaleFactor < 1 || form.ScaleFactor > 20 {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	type ImageCollection struct {
		Upscaled string
		Original string
	}

	var collection []ImageCollection

	for _, i := range form.Images {
		original, config, err := parseImageAndImageConfig(i)

		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}

		upscaled := imaging.Resize(original, config.Width*form.ScaleFactor, config.Height*form.ScaleFactor, imaging.NearestNeighbor)
		upscaledBase64, err := imageEncoders[PNG].encodeImage(upscaled)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}
		originalBase64, err := imageEncoders[PNG].encodeImage(original)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}

		collection = append(collection, ImageCollection{
			Upscaled: upscaledBase64,
			Original: originalBase64,
		})
	}

	rt, _ := template.ParseFiles("result.html")
	err := rt.Execute(w, collection)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}
}
