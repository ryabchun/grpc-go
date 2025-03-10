package main

import (
	"fmt"
	"log"
	"net"

	"grpc-go/configs"
	"grpc-go/internal/db"
	"grpc-go/internal/server"
	"grpc-go/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	cfg := configs.LoadConfig()
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer dbConn.Close()

	grpcServer := grpc.NewServer()
	svc := server.NewService(dbConn)
	proto.RegisterMyServiceServer(grpcServer, svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	log.Printf("gRPC server running on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
