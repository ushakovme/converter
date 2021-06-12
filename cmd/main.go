package main

import (
	"fmt"
	"github.com/ushakovme/converter/pkg/converter/infrastructure"
	"os"
)

func main() {
	fmt.Println("starting converter")
	err := do()
	if err != nil {
		panic(err)
	}
	fmt.Println("converter stopped")
}

func do() error {
	pngFile, err := os.Open("tests/files/bar.png")
	if err != nil {
		return err
	}

	jpgFile, err := os.Create("tests/files/bar.jpg")
	if err != nil {
		return err
	}

	converter := infrastructure.NewConverter()
	return converter.PNGToJPG(pngFile, jpgFile)
}
