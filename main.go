package main

import (
	"net/http"
	"os"

	"github.com/ggicci/httpin"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
)

var (
	r             = http.NewServeMux()
	imageEncoders = make(map[ImageEncoding]ImageStringEncoder)
)

func init() {
	r.HandleFunc("/", index)

	r.Handle("/result",
		alice.New(httpin.NewInput(ResultHandlerForm{})).ThenFunc(result))

	imageEncoders[EncoderPNG] = &PngImageStringEncoder{}
	imageEncoders[EncoderJPEG] = &JpegImageStringEncoder{}
}

func main() {
	err := http.ListenAndServe("0.0.0.0:8080", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		panic(err)
	}
}
