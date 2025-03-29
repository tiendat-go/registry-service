package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	frt "github.com/tiendat-go/common-service/utils/format"
	pb "github.com/tiendat-go/proto-service/gen/registry/v1"
	"github.com/tiendat-go/registry-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	port := frt.GetString(os.Getenv("SERVICE_PORT"), "50051")

	server := service.NewRegistryServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDiscoveryServiceServer(grpcServer, server)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("🛑 Shutting down service registry...")
		grpcServer.GracefulStop()
		os.Exit(0)
	}()

	log.Printf("🔥 Service Registry is running on port:%v", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
