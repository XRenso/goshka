// Package main is the entry point for the API Gateway service.
//
//go:generate swag init -g swagger.go -d .,../internal/handler/rest,../internal/domain -o ../docs --parseDependency --parseInternal
//
// @title           BareRepo API Gateway
// @version         1.0
// @description     API Gateway for GitHub Repository Information Service. Proxies REST requests to the Collector via gRPC.
// @host            localhost:8080
// @BasePath        /

package main
