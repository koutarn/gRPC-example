package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
    "io"
    "errors"

	hellopb "github.com/koutarn/gRPC-example/src/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello %s!", req.GetName()),
	}, nil
}

func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello,%s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
    nameList := make([]string,0)
    for {
        req,err := stream.Recv()

		if errors.Is(err, io.EOF) {
            message := fmt.Sprintf("Hello, %v!", nameList)
            return stream.SendAndClose(&hellopb.HelloResponse{
	            Message: message,
            })
        }

        if err != nil {
            return err
        }

        nameList = append(nameList,req.GetName())
    }
}

// func NewMyServer() *myServer {
// 	return &myServer{}
// }

func main() {

	// リスナーを作成する
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// grpcの鯖のインスタンスを作成する
	s := grpc.NewServer()

	// gRPCサーバーにGreetingServiceを登録
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// サーバーリフレクション設定
	reflection.Register(s)

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
