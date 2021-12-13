package main

import (
	"es/api/handler"
	"es/api/proto/src"
	"es/config"
	"es/internal/repo"
	"es/internal/service"
	"es/pkg/db"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	serviceConfig := config.LoadConfig()
	conn := db.CreateDbConn(serviceConfig)

	estimationRepo := repo.NewEstimationRepo(conn)
	estimationService := service.NewEstimationService(estimationRepo)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", serviceConfig.ServerConfig.Port))
	if err != nil {
		log.Fatalf("faild to listen %v", err)
	}
	Server := grpc.NewServer()
	handlers := handler.NewGrpcHandler(estimationService)
	src.RegisterEsServiceServer(Server, handlers)
	log.Println("starting the grpc server")
	if err := Server.Serve(listen); err != nil {
		log.Fatalf("faild to serve %v", err)
	}

}
