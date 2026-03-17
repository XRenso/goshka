package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/XRenso/barerepo/gateway/docs"
	collectorAdapter "github.com/XRenso/barerepo/gateway/internal/adapter/collector"
	"github.com/XRenso/barerepo/gateway/internal/handler/rest"
	"github.com/XRenso/barerepo/gateway/internal/usecase"
	pb "github.com/XRenso/barerepo/pkg/pb"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "8080"
	}

	conn, err := grpc.NewClient(collectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewRepoServiceClient(conn)
	collClient := collectorAdapter.NewClient(pbClient)
	repoUC := usecase.NewRepoUseCase(collClient)
	handler := rest.NewHandler(repoUC)

	r := gin.Default()
	r.GET("/api/v1/repos/:owner/:repo", handler.GetRepo)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("API Gateway listening on :%s", gatewayPort)
	if err := r.Run(":" + gatewayPort); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
