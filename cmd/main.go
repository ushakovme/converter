package main

import (
	"context"
	"fmt"
	"github.com/ushakovme/converter/pkg/converter/infrastructure"
	proto "github.com/ushakovme/converter/proto/gen/go"
	"google.golang.org/grpc"
	"net"
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
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	converter := infrastructure.NewConverter()
	server := infrastructure.NewServer(converter)
	proto.RegisterConverterServer(grpcServer, server)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	lis, err := net.Listen("tcp", cnf.Port)
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}
