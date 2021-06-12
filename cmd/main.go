package main

import (
	"context"
	"fmt"
	"github.com/ushakovme/converter/pkg/converter/infrastructure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const httpAddr = "127.0.0.1:8000"

func main() {
	fmt.Println("starting converter")
	ctx := getSignalCancelContext()

	if err := serve(ctx); err != nil {
		fmt.Println("converter stopped with error: ")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("converter stopped")
}

func getSignalCancelContext() context.Context {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancelFunc()
	}()
	return ctx
}

func serve(ctx context.Context) error {
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
	server := http.Server{
		Addr: httpAddr,
	}

	server.Handler = mux
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println("HTTP server ListenAndServe: ", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
