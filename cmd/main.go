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

func main() {
	fmt.Println("starting converter")

	err := do()
	if err != nil {
		fmt.Println("converter stopped with error: ")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("converter stopped")
}

func do() error {
	cnf, err := infrastructure.LoadConfig()
	if err != nil {
		return err
	}

	ctx := getSignalCancelContext()

	if err := serve(ctx, cnf); err != nil {
		fmt.Println("converter stopped with error: ")
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
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

func serve(ctx context.Context, cnf infrastructure.Config) error {
	converter := infrastructure.NewConverter()

	mux := http.NewServeMux()
	mux.HandleFunc("/convert", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		pngFile, _, err := request.FormFile("file")
		if err != nil {
			return
		}
		err = converter.PNGToJPG(pngFile, writer)
	})
	server := http.Server{
		Addr: cnf.Port,
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
