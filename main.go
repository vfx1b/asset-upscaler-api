package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	it, _ := template.ParseFiles("index.html")
	err := it.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func result(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal(err)
	}

	file, _, err := r.FormFile("assetFile")
	if err != nil {
		log.Fatal(err)
	}

	i, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(0, 0)
	c, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	i = imaging.Resize(i, c.Width*10, c.Height*10, imaging.NearestNeighbor)

	var buf bytes.Buffer

	png.Encode(&buf, i)

	b64str := base64.StdEncoding.EncodeToString(buf.Bytes())

	fmt.Println(b64str)

	rt, _ := template.ParseFiles("result.html")
	err = rt.Execute(w, b64str)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
