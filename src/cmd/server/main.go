package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "github.com/koutarn/gRPC-example/src/pkg/grpc"
	"google.golang.org/grpc"
)

func main() {

	// リスナーを作成する
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// grpcの鯖を作成する
	s := grpc.NewServer()

	// gRPCサーバーにGreetingServiceを登録
	hellopb.RegisterGreetingServiceServer(s)

	// 鯖を稼動させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// ctrl+cが押されたら停止
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()

}
