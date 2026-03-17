package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	githubAdapter "github.com/XRenso/barerepo/collector/internal/adapter/github"
	grpcHandler "github.com/XRenso/barerepo/collector/internal/handler/grpc"
	"github.com/XRenso/barerepo/collector/internal/usecase"
	pb "github.com/XRenso/barerepo/pkg/pb"
)

func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ghClient := githubAdapter.NewClient()
	repoUC := usecase.NewRepoUseCase(ghClient)
	handler := grpcHandler.NewHandler(repoUC)

	grpcServer := grpc.NewServer()
	pb.RegisterRepoServiceServer(grpcServer, handler)

	log.Printf("Collector gRPC server listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
