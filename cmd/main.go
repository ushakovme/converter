package main

import (
	"fmt"
	"github.com/ushakovme/converter/pkg/converter/infrastructure"
	"net/http"
	"os"
)

const httpAddr = "127.0.0.1:8000"

func main() {
	fmt.Println("starting converter")
	err := do()
	if err != nil {
		panic(err)
	}
	fmt.Println("converter stopped")
}

func do() error {
	converter := infrastructure.NewConverter()

	mux := http.NewServeMux()
	mux.HandleFunc("/convert", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		jpgFile, err := os.Create("tests/files/bar.jpg")
		if err != nil {
			return
		}

		pngFile, _, err := request.FormFile("file")
		if err != nil {
			return
		}
		err = converter.PNGToJPG(pngFile, jpgFile)
	})
	return http.ListenAndServe(httpAddr, mux)
}
